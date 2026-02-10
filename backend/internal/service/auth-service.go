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
		magicLinkTokenPayload, err := service.vTokensRepo.VerifyMagicLink(ctx, hashedToken)
		if err != nil {
			switch {
			case errors.Is(err, errorcodes.ErrNotFound), magicLinkTokenPayload.VerificationType != auth.EmailVerification: // Token not exists or is for another verification
				return errorcodes.ErrInvalid
			default:
				return err
			}
		}

		// Check expiry
		if magicLinkTokenPayload.Exp.Before(time.Now()) {
			return errorcodes.InternalErrExpired
		}

		// Verify user
		err = service.userRepo.Verify(ctx, magicLinkTokenPayload.UserId)
		if err != nil {
			return err
		}

		// Delete token
		_ = service.vTokensRepo.Delete(ctx, magicLinkTokenPayload.VerificationId) // If it fails to delete there are no problems

		return nil
	})
}

/*
Errors
  - ErrInvalid
  - InternalErrExpired
  - Other db errors
*/
func (service *AuthService) VerifyEmailWithOTP(ctx context.Context, verificationId int64, hashedOTP []byte) error {
	return withTransaction(service.db, ctx, func(transaction *sql.Tx) error {
		// Find user related to the token
		// user, err := store.getUserFromEmailVerificationToken(ctx, transaction, plainToken)
		// if err != nil {
		// 	return err
		// }

		// TODO: controlla verificationType, attempts e exp se verifica andata a buon fine

		// TODO: aggiorna attempts se verifica fallita, altrimenti elimina record

		// Update user (email verified)
		// user.IsVerified = true
		// if err := store.setUserIsVerified(ctx, transaction, user.Id); err != nil {
		// 	return err
		// }

		// Clean email verification token
		// if err := store.deleteEmailVerificationToken(ctx, transaction, user.Id); err != nil {
		// 	return err
		// }

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
