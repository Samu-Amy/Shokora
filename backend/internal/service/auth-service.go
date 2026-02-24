package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/errorcodes"
	"github.com/Samu-Amy/Shokora/internal/store"
)

type AuthService struct {
	userRepo           store.UserRepositoryI // TODO: serve (tolta creazione utente)?
	vTokenRepo         store.VTokenRepositoryI
	refreshTokenRepo   store.RefreshTokenRepositoryI
	userSessionRepo    store.UserSessionI
	db                 *sql.DB
	tokenAuthenticator *auth.TokenAuthenticator
}

func NewAuthService(userRepo store.UserRepositoryI, vTokensRepo store.VTokenRepositoryI, refreshTokensRepo store.RefreshTokenRepositoryI, userSessionRepo store.UserSessionI, db *sql.DB, tokenAuthenticator *auth.TokenAuthenticator) *AuthService {
	return &AuthService{userRepo, vTokensRepo, refreshTokensRepo, userSessionRepo, db, tokenAuthenticator}
}

// ------ CREATE USER -----

func (service *AuthService) CreateUser(ctx context.Context, user *store.User) error {
	// TODO: fare transaction per creazione user, stats and settings (oppure crearle nell'update se non esistono)?

	return app.store.User.Create(ctx, user)
}

// ----- CREATE VERIFICATION TOKENS -----

func (service *AuthService) CreateVerificationTokensWithRetries(ctx context.Context, user *store.User, verificationTokens *auth.VerificationTokens) (*int64, error) {

	ctxWithTimeout, cancel := context.WithTimeout(ctx, REGENERATE_TOKEN_TIMEOUT)
	defer cancel()

	// Create Tokens in db
	verificationId, err := service.vTokenRepo.CreateTokens(ctxWithTimeout, user.Id, verificationTokens)
	if err == nil {
		return verificationId, nil // OK, return no error

	} else if !errors.Is(err, errorcodes.InternalErrDuplicateToken) {
		return nil, err // Error (can't retry)
	}

	// Retries
	for range service.tokenAuthenticator.MaxRetries - 1 {

		// Timeout
		select {
		case <-ctxWithTimeout.Done():
			return nil, ctxWithTimeout.Err()
		default:
		}

		// Regenerate Tokens (if error is "duplicate token")
		switch {

		// Duplicated Magic Link Token
		case errors.Is(err, errorcodes.InternalErrDuplicateToken):
			err = service.tokenAuthenticator.RegenerateMagicLinkToken(verificationTokens)
			if err != nil {
				continue // skip iteration
			}

			verificationId, err = service.vTokenRepo.CreateTokens(ctx, user.Id, verificationTokens)
			if err == nil {
				return verificationId, nil // OK, return no error
			}

		default:
			return nil, err // Error is not solvable (not "duplicate token") -> return it
		}
	}

	return nil, errorcodes.ErrMaxRetriesExceeded // Couldn't regenerate and save token successfully -> return error "max retries exceeded"
}

// ----- VERIFY EMAIL  -----

/*
Errors
  - ErrInvalid
  - Other db errors
*/
func (service *AuthService) VerifyEmailWithToken(ctx context.Context, hashedToken []byte) error {

	// Verify and Get data
	magicLinkTokenQueryData, err := service.vTokenRepo.VerifyMagicLink(ctx, hashedToken, auth.EmailVerification)
	if err != nil {
		// log.Printf("Verify OTP Error: %v", err)
		switch {
		case errors.Is(err, errorcodes.ErrNotFound): // Token not valid
			return errorcodes.ErrInvalid
		default:
			return err
		}
	}

	// Verify user
	err = service.userRepo.Verify(ctx, magicLinkTokenQueryData.UserId)
	if err != nil {
		// log.Printf("Verify User Error: %v", err)
		return err
	}

	// Delete token
	_ = service.vTokenRepo.Delete(ctx, magicLinkTokenQueryData.VerificationId) // If it fails to delete there are no problems

	return nil
}

/*
Errors
  - ErrInvalid
  - InternalErrExpired
  - ErrMaxAttemptsExceeded
  - Other db errors
*/
func (service *AuthService) VerifyEmailWithOTP(ctx context.Context, verificationId int64, hashedOTP []byte, maxAttempts uint8) error {

	// Get data
	otpQueryData, err := service.verifyOtp(ctx, verificationId, hashedOTP, maxAttempts, auth.EmailVerification)
	if err != nil {
		// log.Printf("Verify OTP Error: %v", err)
		return err
	}

	// Verify user
	err = service.userRepo.Verify(ctx, otpQueryData.UserId)
	if err != nil {
		// log.Printf("Verify User Error: %v", err)
		return err
	}

	// Delete token
	_ = service.vTokenRepo.Delete(ctx, verificationId) // If it fails to delete there are no problems

	return nil
}

// ----- RESET PASSWORD  -----

// ----- TWO FACTOR AUTH  -----

// ----- CREATE AND ROTATE REFRESH TOKENS -----

// Create
func (service *AuthService) CreateRefreshToken(ctx context.Context, session *store.UserSession, refreshToken *store.RefreshToken, sessionExp, tokenExp time.Duration) error {
	// TODO: fai transaction per creare sia sessione che refresh token

	return withTransaction(service.db, ctx, func(tx *sql.Tx) error {
		// Create session
		err := service.userSessionRepo.Create(ctx, tx, session, sessionExp)
		if err != nil {
			return err
		}

		// Create token
		return service.refreshTokenRepo.CreateToken(ctx, tx, refreshToken, tokenExp)
	})

}

// Rotate
func (service *AuthService) RotateRefreshToken(ctx context.Context, oldHashedToken []byte, newRefreshToken *store.RefreshToken) error { // TODO: ritorna dati (es. expires_at per cookies)
	err := withTransaction(service.db, ctx, func(tx *sql.Tx) error {

		// Get token
		oldRefreshToken, err := service.refreshTokenRepo.GetToken(ctx, tx, oldHashedToken)
		if err != nil {
			return err
		}

		// Validate token - Expired or revoked
		if oldRefreshToken.ExpiresAt.Before(time.Now()) || oldRefreshToken.RevokedAt != nil {
			return errorcodes.ErrInvalid
		}

		// Validate token - wrong RefreshToken, User or Session id
		if oldRefreshToken.UserId != newRefreshToken.UserId || oldRefreshToken.SessionId != newRefreshToken.SessionId || oldRefreshToken.Id != newRefreshToken.Replaces {
			return errorcodes.ErrInvalid
		}

		// TODO: implementa estensione expiry - usando costanti "SESSION_EXTENSION_<...>" (e nel token nuovo usa il conto di estensioni da quello vecchio (più eventualmente quella appena fatta))

		// Create new token
		// (create token using same session_id of the old one and using its id as replaces)
		err = service.refreshTokenRepo.CreateToken(ctx, tx, newRefreshToken)
		if err != nil {
			return err
		}

		if oldRefreshToken.Id == nil || newRefreshToken.CreatedAt == nil {
			return errorcodes.ErrInvalid
		}

		// Revoke (update) old token
		err = service.refreshTokenRepo.RevokeTokenById(ctx, tx, *oldRefreshToken.Id, *newRefreshToken.CreatedAt)
		if err != nil {
			// TODO: gestire token not found or already revoked (?)
			switch {
			case errors.Is(err, errorcodes.InternalErrNoRowsAffected):
				return errorcodes.InternalErrTokenNotFoundOrAlreadyRevoked
			default:
				return err
			}
		}

		return nil
	})

	// Revoke session (in case of Reuse Detection)
	if errors.Is(err, errorcodes.InternalErrReusedToken) {
		err = service.refreshTokenRepo.DeleteSessionById(ctx, newRefreshToken.UserId, newRefreshToken.SessionId)
	}

	return err
}
