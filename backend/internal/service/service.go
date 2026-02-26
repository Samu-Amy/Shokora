package service

import (
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/auth"
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

func NewService(db *sql.DB, mailer mailer.IClient, store *store.Storage, jwtAuthenticator *auth.JWTAuthenticator, tokenAuthenticator *auth.TokenAuthenticator, authServiceConfig authservice.AuthServiceConfig) *Service {
	return &Service{
		Auth: authservice.NewService(db, mailer, store.User, store.VToken, store.RefreshToken, store.UserSession, jwtAuthenticator, tokenAuthenticator, authServiceConfig),
		User: userservice.NewService(db, store.User),
	}
}
