package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/store"
)

type AuthService struct {
	userRepo           store.UserRepositoryI
	vTokensRepo        store.VTokensRepositoryI
	db                 *sql.DB
	tokenAuthenticator *auth.TokenAuthenticator
}

func NewAuthService(userRepo store.UserRepositoryI, vTokensRepo store.VTokensRepositoryI, db *sql.DB, tokenAuthenticator *auth.TokenAuthenticator) *AuthService {
	return &AuthService{userRepo, vTokensRepo, db, tokenAuthenticator}
}

// ----- CREATE USER -----

func (service *AuthService) CreateUserAndEmailVerificationTokensWithRetries(ctx context.Context, user *store.User, verificationTokens *auth.VerificationTokens) error {
	// Create user
	if err := service.userRepo.Create(ctx, user); err != nil {
		return err
	}

	// TODO: fare transaction per creazione user, stats and settings

	ctxWithTimeout, cancel := context.WithTimeout(ctx, regenerate_token_timeout)
	defer cancel()

	// Create verification tokens
	return service.createVerificationTokensWithRetries(ctxWithTimeout, user.Id, verificationTokens)
}

// ----- VERIFY EMAIL  -----

func (service *AuthService) VerifyEmail(ctx context.Context, plainToken string) error { // TODO: passare plain token e verificare con funzione util (?)
	return withTransaction(service.db, ctx, func(transaction *sql.Tx) error {
		// // Find user related to the token
		// user, err := store.getUserFromEmailVerificationToken(ctx, transaction, plainToken)
		// if err != nil {
		// 	return err
		// }

		// // Update user (email verified)
		// user.IsVerified = true
		// if err := store.setUserIsVerified(ctx, transaction, user.Id); err != nil {
		// 	return err
		// }

		// // Clean email verification token
		// if err := store.deleteEmailVerificationToken(ctx, transaction, user.Id); err != nil {
		// 	return err
		// }

		return nil
	})
}

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

// ----- UTILS -----

func (service *AuthService) createVerificationTokensWithRetries(ctx context.Context, userId int64, verificationTokens *auth.VerificationTokens) error {
	// Create Tokens in db
	err := service.vTokensRepo.CreateTokens(ctx, userId, verificationTokens)
	if err == nil {
		return nil // OK, return no error
	}

	// Retry (regenerate tokens)
	for range service.tokenAuthenticator.MaxRetries - 1 {

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Regenerate Tokens (if error is "duplicate token")
		switch {

		// Magic Link Token
		case errors.Is(err, store.ErrDuplicateToken):
			err2 := service.tokenAuthenticator.RegenerateMagicLinkToken(verificationTokens)
			if err2 != nil {
				continue // skip iteration
			}

			err = service.vTokensRepo.UpdateMagicLinkToken(ctx, userId, verificationTokens.VerificationType, verificationTokens.HashedMagicLinkToken, verificationTokens.MagicLinkTokenExp)
			if err == nil {
				return nil // OK, return no error
			}

		// 	// OTP
		// TODO: implementare controllo duplicazione e rigenerazione OTP?
		// case errors.Is(err, store.ErrDuplicateOTP):
		// 	err2 := service.tokenAuthenticator.RegenerateOTP(verificationTokens)

		// 	// TODO: fai query per aggiornare otp
		default:
			return err // Error is not solvable (not "duplicate token") -> return it
		}
	}

	return ErrMaxRetriesExceeded // Couldn't regenerate and save token successfully -> return error "max retries exceeded"
}
