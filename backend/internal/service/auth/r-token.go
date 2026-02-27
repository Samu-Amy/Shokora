package authservice

import (
	"context"
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/auth"
	rtoken "github.com/Samu-Amy/Shokora/internal/store/refresh-token.go"
)

// ----- CREATE AND ROTATE REFRESH TOKENS -----

// Create
func (service *AuthService) createNewRefreshToken(ctx context.Context, userId int64) (*rtoken.CreateRefreshTokenDto, error) {
	// TODO: genera refresh token qua e ritornalo dopo averlo messo nel db

	// Create session
	// Generate Token
	// Create Refresh Token
	var createRefreshTokenDto = &rtoken.CreateRefreshTokenDto{}

	err := service.txManager.WithTx(ctx, func(tx *sql.Tx) error {
		// Create session
		sessionId, err := service.userSessionRepo.Create(ctx, tx, userId, service.config.Token.SessionMaxExp)
		if err != nil {
			return err
		}

		// Generate refresh token
		plainRefreshToken, refreshToken, err := service.generateRefreshToken(sessionId)
		if err != nil {
			return err
		}

		createRefreshTokenDto.PlainToken = *plainRefreshToken

		// Create refresh token in db
		err = service.refreshTokenRepo.CreateToken(ctx, tx, refreshToken, tokenExp)
		if err != nil {
			return err
		}

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

/*
Generate Refresh Token

Return:
  - plainToken
  - refreshToken
  - error
*/
func (service *AuthService) generateRefreshToken(sessionId int64) (*string, *rtoken.RefreshToken, error) {
	plainToken, err := auth.GenerateBase64Token(service.config.Token.RefreshTokenByteSize)
	if err != nil {
		return nil, nil, err
	}

	// Hash token
	hashedToken := auth.HashBase64Token(plainToken)

	refreshToken := &rtoken.RefreshToken{
		SessionId: sessionId,
		TokenHash: hashedToken,
	}
	return plainToken, refreshToken, nil
}
