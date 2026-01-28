package service

import (
	"context"
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/store"
)

type AuthService struct {
	User    store.UserRepository
	VTokens store.VTokensRepository
	db      *sql.DB
}

func NewAuthService(User store.UserRepository, VTokens store.VTokensRepository, db *sql.DB) *AuthService {
	return &AuthService{User, VTokens, db}
}

// type AuthService interface {
// 	CreateUserAndEmailVerificationToken(ctx context.Context, user *store.User, verificationTokens *auth.VerificationTokens) error
// }

func (service *AuthService) CreateUserAndEmailVerificationToken(ctx context.Context, user *store.User, verificationTokens *auth.VerificationTokens) error {
	return withTransaction(service.db, ctx, func(transaction *sql.Tx) error {
		// Create user
		if err := service.User.Create(ctx, transaction, user); err != nil {
			return err
		}

		// Create verification
		if err := service.VTokens.CreateTokens(ctx, transaction, verificationTokens, user.Id); err != nil {
			return err
		}

		return nil
	})
}
