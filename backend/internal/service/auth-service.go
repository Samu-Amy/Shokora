package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/errorcodes"
	"github.com/Samu-Amy/Shokora/internal/store"
)

type AuthService struct {
	userRepo           store.UserRepositoryI // TODO: serve (tolta creazione utente)?
	vTokensRepo        store.VTokensRepositoryI
	db                 *sql.DB
	tokenAuthenticator *auth.TokenAuthenticator
}

func NewAuthService(userRepo store.UserRepositoryI, vTokensRepo store.VTokensRepositoryI, db *sql.DB, tokenAuthenticator *auth.TokenAuthenticator) *AuthService {
	return &AuthService{userRepo, vTokensRepo, db, tokenAuthenticator}
}

// ----- CREATE TOKENS -----

func (service *AuthService) CreateVerificationTokensWithRetries(ctx context.Context, user *store.User, verificationTokens *auth.VerificationTokens) (*int64, error) {

	ctxWithTimeout, cancel := context.WithTimeout(ctx, regenerate_token_timeout)
	defer cancel()

	// Create Tokens in db
	verificationId, err := service.vTokensRepo.CreateTokens(ctxWithTimeout, user.Id, verificationTokens)
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

			verificationId, err = service.vTokensRepo.CreateTokens(ctx, user.Id, verificationTokens)
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
	magicLinkTokenQueryData, err := service.vTokensRepo.VerifyMagicLink(ctx, hashedToken, auth.EmailVerification)
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
	_ = service.vTokensRepo.Delete(ctx, magicLinkTokenQueryData.VerificationId) // If it fails to delete there are no problems

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
	_ = service.vTokensRepo.Delete(ctx, verificationId) // If it fails to delete there are no problems

	return nil
}

// ----- RESET PASSWORD  -----

// ----- TWO FACTOR AUTH  -----
