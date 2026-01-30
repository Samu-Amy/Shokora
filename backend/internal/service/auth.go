package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func (service *AuthService) CreateUserAndEmailVerificationTokens(ctx context.Context, user *store.User, verificationTokens *auth.VerificationTokens) error {
	// Create user
	if err := service.userRepo.Create(ctx, user); err != nil {
		return err
	}

	// TODO: fare transaction per creazione user, stats and settings

	ctx, cancel := context.WithTimeout(ctx, regenerate_token_timeout)
	defer cancel()

	// Create verification tokens
	return service.createVerificationTokensWithRetries(ctx, user.Id, verificationTokens)
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
	err := service.vTokensRepo.CreateTokens(ctx, userId, verificationTokens)
	if err == nil {
		return nil
	}

	fmt.Print("\n\n Creato token \n\n")

	for range service.tokenAuthenticator.MaxRetries {

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Regenerate Tokens
		switch {

		// Magic Link Token
		case errors.Is(err, store.ErrDuplicateToken):
			err = service.tokenAuthenticator.RegenerateMagicLinkToken(verificationTokens)

			// TODO: fai query per aggiornare token
			// OTP
		case errors.Is(err, store.ErrDuplicateOTP):
			err = service.tokenAuthenticator.RegenerateOTP(verificationTokens)

			// TODO: fai query per aggiornare otp
		default:
			return err
		}

		continue
	}

	return ErrMaxRetriesExceeded
}
