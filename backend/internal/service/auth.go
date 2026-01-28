package service

import (
	"context"
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/store"
)

type AuthService struct {
	userRepo           store.UserRepositoryI
	vTokensRepo        store.VTokensRepositoryI
	db                 *sql.DB
	tokenAuthenticator auth.TokenAuthenticatorI
}

func NewAuthService(userRepo store.UserRepositoryI, vTokensRepo store.VTokensRepositoryI, db *sql.DB, tokenAuthenticator auth.TokenAuthenticatorI) *AuthService {
	return &AuthService{userRepo, vTokensRepo, db, tokenAuthenticator}
}

// ----- CREATE -----

func (service *AuthService) CreateUserAndEmailVerificationTokens(ctx context.Context, user *store.User, verificationTokens *auth.VerificationTokens) error {
	return withTransaction(service.db, ctx, func(transaction *sql.Tx) error {
		// Create user
		if err := service.userRepo.Create(ctx, transaction, user); err != nil {
			return err
		}

		// Create verification
		if err := service.vTokensRepo.CreateTokens(ctx, transaction, verificationTokens, user.Id); err != nil {
			return err // TODO: fai retries (usando tokenAuthenticator per rigenerare i token)
		}

		return nil
	})
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
