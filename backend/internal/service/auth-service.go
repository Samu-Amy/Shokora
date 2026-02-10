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
  - InternalErrExpired
  - Other db errors
*/
func (service *AuthService) VerifyEmailWithToken(ctx context.Context, hashedToken []byte) error {
	return withTransaction(service.db, ctx, func(transaction *sql.Tx) error {

		// Get data
		magicLinkTokenQueryData, err := service.vTokensRepo.GetMagicLinkData(ctx, hashedToken)
		if err != nil {
			switch {
			case errors.Is(err, errorcodes.ErrNotFound), magicLinkTokenQueryData.VerificationType != auth.EmailVerification: // Token not exists or is for another verification
				return errorcodes.ErrInvalid
			default:
				return err
			}
		}

		// Check expiry
		if magicLinkTokenQueryData.Exp.Before(time.Now()) {
			return errorcodes.InternalErrExpired
		}

		// Verify user
		err = service.userRepo.Verify(ctx, magicLinkTokenQueryData.UserId)
		if err != nil {
			return err
		}

		// Delete token
		_ = service.vTokensRepo.Delete(ctx, magicLinkTokenQueryData.VerificationId) // If it fails to delete there are no problems

		return nil
	})
}

/*
Errors
  - ErrInvalid
  - InternalErrExpired
  - ErrMaxAttemptsExceeded
  - Other db errors
*/
func (service *AuthService) VerifyEmailWithOTP(ctx context.Context, verificationId int64, hashedOTP []byte, maxAttempts uint8) error {
	return withTransaction(service.db, ctx, func(transaction *sql.Tx) error {

		var verificationErr error

		// Get data
		otpQueryData, err := service.vTokensRepo.GetOTPData(ctx, verificationId, hashedOTP)
		if err != nil {
			switch {
			case errors.Is(err, errorcodes.ErrNotFound): // OTP Not valid
				return errorcodes.ErrInvalid
			default:
				return err // db/query error
			}
		}

		// Check expiry
		if otpQueryData.Exp.Before(time.Now()) {

			verificationErr = errorcodes.InternalErrExpired

		} else if otpQueryData.Attempts >= maxAttempts {

			// Check attempts
			verificationErr = errorcodes.ErrMaxAttemptsExceeded
		}

		// Increment attempts and Handle errors
		if verificationErr != nil {
			err := service.vTokensRepo.UpdateAttempts(ctx, verificationId, maxAttempts)
			if err != nil {
				verificationErr = err // MaxAttemptsExceeded or db/query error
			}

			return verificationErr
		}

		// Verify user
		err = service.userRepo.Verify(ctx, otpQueryData.UserId)
		if err != nil {
			return err
		}

		// Delete token
		_ = service.vTokensRepo.Delete(ctx, verificationId) // If it fails to delete there are no problems

		return nil
	})
}

// ----- RESET PASSWORD  -----

// ----- TWO FACTOR AUTH  -----

// ----- DELETE USER -----

func (service *AuthService) DeleteUserAndEmailVerificationToken(ctx context.Context, userId int64) error {
	return withTransaction(service.db, ctx, func(transaction *sql.Tx) error {
		if err := service.userRepo.Delete(ctx, transaction, userId); err != nil {
			return err
		}

		// if err := service.vTokensRepo.deleteEmailVerificationToken(ctx, transaction, userId); err != nil {
		// 	return err
		// }

		return nil
	})
}
