package authservice

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
		plainRefreshToken, refreshToken, err := service.createRefreshTokenWithRetries(ctx, tx, sessionId, nil, nil)
		if err != nil {
			return err
		}

		// Update refresh token dto
		createRefreshTokenDto.PlainRefreshToken = plainRefreshToken
		createRefreshTokenDto.RefreshTokenExpiresAt = refreshToken.ExpiresAt

		return nil
	})

	if err != nil {
		return nil, err
	}

	return createRefreshTokenDto, nil
}

/*
Rotate Refresh Token

Returns:
  - *payloads.AuthTokensDto with PlainRefreshToken and RefreshTokenExpiresAt
  - int64: userId
  - error: interrors (to be parsed into domerrors)
*/
func (service *AuthService) rotateRefreshToken(ctx context.Context, oldHashedToken []byte) (*payloads.AuthTokensDto, int64, error) {

	var createRefreshTokenDto = &payloads.AuthTokensDto{}
	var sessionId int64 = -1
	var userId int64 = -1

	err := service.txManager.WithTx(ctx, func(tx *sql.Tx) error {

		// Get old refresh token and session data
		oldTokenAndSessionData, err := service.refreshTokenRepo.GetByToken(ctx, tx, oldHashedToken)
		if err != nil {
			service.logger.Warnw("Error getting the old refresh token", "error", err)
			return err
		}

		// Set data outside transaction
		sessionId = oldTokenAndSessionData.SessionId
		userId = oldTokenAndSessionData.UserId

		// Validate token - Expired
		if oldTokenAndSessionData.ExpiresAt.Before(time.Now()) {
			service.logger.Warn("Old refresh token expired")
			return interrors.IErrExpired
		}

		// Validate token - Revoked (there is already a token that replaces it)
		if oldTokenAndSessionData.RevokedAt != nil {
			service.logger.Warn("Old refresh token expired")
			return interrors.IErrReusedToken
		}

		// Try to extend expiration
		var newTokenExpiresAt time.Time

		if time.Until(oldTokenAndSessionData.ExpiresAt) <= auth.SessionExtensionCondition {
			newExpiresAt := oldTokenAndSessionData.ExpiresAt.Add(auth.SessionExtensionDuration)

			// Extend expiration by min(oldToken ExpiresAt + SessionExtensionDuration, sessionExpiresAt)
			if newExpiresAt.Before(oldTokenAndSessionData.SessionExpiresAt) {
				newTokenExpiresAt = newExpiresAt
			} else {
				newTokenExpiresAt = oldTokenAndSessionData.SessionExpiresAt
			}
		} else {
			newTokenExpiresAt = oldTokenAndSessionData.ExpiresAt
		}

		// Create Refresh Token in db
		newPlainRefreshToken, newRefreshToken, err := service.createRefreshTokenWithRetries(ctx, tx, oldTokenAndSessionData.SessionId, &oldTokenAndSessionData.Id, &newTokenExpiresAt)
		if err != nil {
			return err
		}

		// Update refresh token dto
		createRefreshTokenDto.PlainRefreshToken = newPlainRefreshToken
		createRefreshTokenDto.RefreshTokenExpiresAt = newTokenExpiresAt

		// Validate data for revoked_at
		newTokenCreatedAt := newRefreshToken.CreatedAt
		if newTokenCreatedAt.IsZero() {
			newTokenCreatedAt = time.Now()
		}

		// Revoke (update) old token
		err = service.refreshTokenRepo.RevokeById(ctx, tx, oldTokenAndSessionData.Id, newTokenCreatedAt)
		if err != nil {
			service.logger.Warnw("Error updating the old refresh token", "error", err)
			return err
		}

		return nil
	})

	// Revoke session (in case of Reuse Detection)
	if errors.Is(err, interrors.IErrReusedToken) {
		service.logger.Warnw("Reused Refresh Token", "error", err)

		if sessionId != -1 {
			err = service.userSessionRepo.Delete(ctx, sessionId) // TODO: controlla che sessionId sia sempre disponibile in questo caso (GetByToken non dovrebbe mai ritornare IErrReusedToken, quindi dovrebbe essere ok)
			if err != nil {
				service.logger.Warnw("Error deleting Session", "error", err)
			}
		}
	}

	return createRefreshTokenDto, userId, err
}

// ----- GENERATE / CREATE TOKEN -----

/*
Generate Refresh Token

Return:
  - plainToken
  - refreshToken (rtoken.RefreshToken)
  - error
*/
func (service *AuthService) generateRefreshToken(sessionId int64, replaces *int64, expiresAt *time.Time) (string, *rtoken.RefreshToken, error) {
	plainToken, err := auth.GenerateBase64Token(service.config.Token.RefreshTokenByteSize)
	if err != nil {
		return plainToken, nil, err
	}

	// Hash token
	hashedToken := auth.HashBase64Token(plainToken)

	// ExpiresAt
	tokenExpiresAt := time.Now().Add(service.config.Token.RefreshTokenExp)
	if expiresAt != nil {
		tokenExpiresAt = *expiresAt
	}

	refreshToken := &rtoken.RefreshToken{
		SessionId: sessionId,
		TokenHash: hashedToken,
		Replaces:  replaces,
		ExpiresAt: tokenExpiresAt,
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
func (service *AuthService) createRefreshTokenWithRetries(ctx context.Context, tx *sql.Tx, sessionId int64, replaces *int64, expiresAt *time.Time) (string, *rtoken.RefreshToken, error) {
	// Generate token
	plainRefreshToken, refreshToken, err := service.generateRefreshToken(sessionId, replaces, expiresAt)
	if err != nil {
		service.logger.Warnw("Error generating refresh token", "error", err)
		return "", nil, err
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, auth.RegenerateTokenTimeout)
	defer cancel()

	// Create token in db
	err = service.refreshTokenRepo.Create(ctxWithTimeout, tx, refreshToken)
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
			plainRefreshToken, refreshToken, err = service.generateRefreshToken(sessionId, replaces, expiresAt)
			if err != nil {
				service.logger.Warnw("Error rigenerating refresh token", "error", err)

				// skip iteration (reset err to IErrDuplicateToken)
				err = interrors.IErrDuplicateToken
				continue
			}

			err = service.refreshTokenRepo.Create(ctxWithTimeout, tx, refreshToken)
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
