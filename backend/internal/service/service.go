package service

import (
	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/database"
	"github.com/Samu-Amy/Shokora/internal/mailer"
	authservice "github.com/Samu-Amy/Shokora/internal/service/auth"
	userservice "github.com/Samu-Amy/Shokora/internal/service/user"
	"github.com/Samu-Amy/Shokora/internal/store"
)

type Service struct {
	Auth *authservice.AuthService // Verification Tokens, Refresh Tokens, Permissions
	User *userservice.UserService // (settings, stats, achievements)
	// Menu (menu sections -> productsId)
	// Shop
	// Orders
}

func NewService(txManager database.ITransactionManager, store *store.Storage, mailer mailer.IClient, jwtAuthenticator *auth.JWTAuthenticator, tokenAuthenticator *auth.TokenAuthenticator, authServiceConfig authservice.AuthServiceConfig) *Service {
	return &Service{
		Auth: authservice.NewService(txManager, store.User, store.VToken, store.RefreshToken, store.UserSession, mailer, jwtAuthenticator, tokenAuthenticator, authServiceConfig),
		User: userservice.NewService(txManager, store.User),
	}
}
