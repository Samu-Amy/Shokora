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
	// TODO: togliere transaction (l'utente deve essere comunque creato, se non dovesse riuscire a creare i token si possono rigenerare)
	// return withTransaction(service.db, ctx, func(transaction *sql.Tx) error {
	// Create user
	if err := service.userRepo.Create(ctx, user); err != nil {
		return err
	}

	// Create verification tokens
	for range service.tokenAuthenticator.MaxRetries {

		// TODO: (aggiungere timeout - magari nell'handler?), utile per retry e altre operazioni che prendono tempo, va bene?
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := service.vTokensRepo.CreateTokens(ctx, verificationTokens, user.Id); err != nil {
			switch {
			case errors.Is(err, store.ErrDuplicateMagicLinkToken):

				// Regenerate email verification token
				newMagicLinkToken, err2 := service.tokenAuthenticator.GenerateMagicLinkToken()
				if err2 != nil {
					continue
				}

				verificationTokens.PlainMagicLinkToken = newMagicLinkToken
				verificationTokens.HashedMagicLinkToken = service.tokenAuthenticator.HashMagicLinkToken(newMagicLinkToken)

			case errors.Is(err, store.ErrDuplicateOTP):

				// Regenerate otp
				newOTP, err2 := service.tokenAuthenticator.GenerateOTP()
				if err2 != nil {
					continue
				}

				verificationTokens.PlainOTP = newOTP
				verificationTokens.HashedOTP = service.tokenAuthenticator.HashOTP(newOTP, verificationTokens.VerificationType)
			default:
				return err
			}

			continue
		}

		return nil
	}

	return fmt.Errorf("max_retries") // TODO: crea errore apposta (?)
	// })
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
