package authservice

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/auth"
	interrors "github.com/Samu-Amy/Shokora/internal/errors/int"
	rtoken "github.com/Samu-Amy/Shokora/internal/store/refresh-token.go"
)

// ----- CREATE AND ROTATE REFRESH TOKENS -----

// Create
func (service *AuthService) createNewSessionAndRefreshToken(ctx context.Context, userId int64) (*payloads.AuthTokensDto, error) {

	var createRefreshTokenDto = &payloads.AuthTokensDto{}

	err := service.txManager.WithTx(ctx, func(tx *sql.Tx) error {

		// Create session
		sessionId, err := service.userSessionRepo.Create(ctx, tx, userId, service.config.Token.SessionExp)
		if err != nil {
			service.logger.Warnw("Error creating session in db", "error", err)
			return err
		}

		// Create Refresh Token in db
		plainRefreshToken, refreshToken, err := service.createRefreshTokenWithRetries(ctx, tx, sessionId, nil)
		if err != nil {
			return err
		}

		createRefreshTokenDto.PlainRefreshToken = plainRefreshToken

		createRefreshTokenDto.RefreshTokenExpiresAt = refreshToken.ExpiresAt

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createRefreshTokenDto, nil
}

// TODO: nel login fai anche delete di tutti i refresh token scaduti per quell'utente (o in generale?) - ottenere un l'ultimo token creato per ogni session_id (join con order by) e se è scaduto -> sessione scaduta (?)
// TODO: aggiorna a session + refresh token

// Rotate
// func (service *AuthService) rotateRefreshToken(ctx context.Context, oldHashedToken []byte, newRefreshToken *rtoken.RefreshToken) error {
// 	err := db.WithTransaction(service.db, ctx, func(tx *sql.Tx) error {

// 		// Get token
// 		oldRefreshToken, err := service.refreshTokenRepo.GetToken(ctx, tx, oldHashedToken)
// 		if err != nil {
// 			return err
// 		}

// 		// Validate token - Expired or revoked
// 		if oldRefreshToken.ExpiresAt.Before(time.Now()) || oldRefreshToken.RevokedAt != nil {
// 			return errorcodes.ErrInvalid
// 		}

// 		// Validate token - wrong RefreshToken, User or Session id
// 		if oldRefreshToken.UserId != newRefreshToken.UserId || oldRefreshToken.SessionId != newRefreshToken.SessionId || oldRefreshToken.Id != newRefreshToken.Replaces {
// 			return errorcodes.ErrInvalid
// 		}

// 		// TODO: implementa estensione expiry - usando costanti "SessionExtensionDuration" e "SessionExtensionCondition" (e nel token nuovo usa il conto di estensioni da quello vecchio (più eventualmente quella appena fatta))

// 		// Create new token
// 		// (create token using same session_id of the old one and using its id as replaces)
// 		err = service.refreshTokenRepo.CreateToken(ctx, tx, newRefreshToken)
// 		if err != nil {
// 			return err
// 		}

// 		if oldRefreshToken.Id == nil || newRefreshToken.CreatedAt == nil {
// 			return errorcodes.ErrInvalid
// 		}

// 		// Revoke (update) old token
// 		err = service.refreshTokenRepo.RevokeTokenById(ctx, tx, *oldRefreshToken.Id, *newRefreshToken.CreatedAt)
// 		if err != nil {
// 			// TODO: gestire token not found or already revoked (?)
// 			switch {
// 			case errors.Is(err, errorcodes.InternalErrNoRowsAffected):
// 				return errorcodes.InternalErrTokenNotFoundOrAlreadyRevoked
// 			default:
// 				return err
// 			}
// 		}

// 		return nil
// 	})

// 	// Revoke session (in case of Reuse Detection)
// 	if errors.Is(err, errorcodes.InternalErrReusedToken) {
// 		err = service.userSessionRepo.DeleteSessionById(ctx, newRefreshToken.UserId, newRefreshToken.SessionId)
// 	}

// 	return err
// }

// ----- GENERATE / CREATE TOKEN -----

/*
Generate Refresh Token

Return:
  - plainToken
  - refreshToken (rtoken.RefreshToken)
  - error
*/
func (service *AuthService) generateRefreshToken(sessionId int64, replaces *int64) (string, *rtoken.RefreshToken, error) {
	plainToken, err := auth.GenerateBase64Token(service.config.Token.RefreshTokenByteSize)
	if err != nil {
		return plainToken, nil, err
	}

	// Hash token
	hashedToken := auth.HashBase64Token(plainToken)

	refreshToken := &rtoken.RefreshToken{
		SessionId: sessionId,
		TokenHash: hashedToken,
		Replaces:  replaces,
	}

	return plainToken, refreshToken, nil
}

/*
Create Refresh Token

Return:
  - plainToken
  - refreshToken (rtoken.RefreshToken)
  - error
*/
func (service *AuthService) createRefreshTokenWithRetries(ctx context.Context, tx *sql.Tx, sessionId int64, replaces *int64) (string, *rtoken.RefreshToken, error) {
	// Generate token
	plainRefreshToken, refreshToken, err := service.generateRefreshToken(sessionId, replaces)
	if err != nil {
		service.logger.Warnw("Error generating refresh token", "error", err)
		return "", nil, err
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, auth.RegenerateTokenTimeout)
	defer cancel()

	// Create token in db
	err = service.refreshTokenRepo.Create(ctxWithTimeout, tx, refreshToken, service.config.Token.RefreshTokenExp)
	if err == nil {
		return plainRefreshToken, refreshToken, nil // Token Created (OK)

	} else if !errors.Is(err, interrors.IErrDuplicateToken) {
		service.logger.Warnw("Error creating refresh token in db", "error", err)
		return "", nil, err // Error (can't retry)
	}

	// Retries (it's MaxRetries - 1 because one try is already done)
	for range service.tokenAuthenticator.MaxRetries - 1 {

		// Timeout
		select {
		case <-ctxWithTimeout.Done():
			return "", nil, ctxWithTimeout.Err()
		default:
		}

		// Regenerate Token (if error is "duplicate token")
		switch {

		case errors.Is(err, interrors.IErrDuplicateToken):
			plainRefreshToken, refreshToken, err = service.generateRefreshToken(sessionId, replaces)
			if err != nil {
				service.logger.Warnw("Error rigenerating refresh token", "error", err)

				// skip iteration (reset err to IErrDuplicateToken)
				err = interrors.IErrDuplicateToken
				continue
			}

			err = service.refreshTokenRepo.Create(ctxWithTimeout, tx, refreshToken, service.config.Token.RefreshTokenExp)
			if err == nil {
				return plainRefreshToken, refreshToken, nil // Tokens Created (OK)
			}

		default:
			service.logger.Warnw("Error during refresh token creation retries", "error", err)
			return "", nil, err // Error is not solvable (not "duplicate token") -> return it
		}
	}

	return "", nil, interrors.IErrMaxRetriesExceeded // Couldn't regenerate and save token successfully -> return error "max retries exceeded"
}
