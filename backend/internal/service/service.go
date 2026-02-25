package service

import (
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/auth"
	authservice "github.com/Samu-Amy/Shokora/internal/service/auth"
	"github.com/Samu-Amy/Shokora/internal/store"
)

type Service struct {
	Auth *authservice.AuthService // Verification Tokens, Refresh Tokens, Permissions
	// Users (settings, stats, achievements)
	// Menu (menu sections -> productsId)
	// Shop
	// Orders
}

func NewService(db *sql.DB, store *store.Storage, tokenAuthenticator *auth.TokenAuthenticator) *Service {
	return &Service{
		Auth: authservice.NewService(store.User, store.VToken, store.RefreshToken, store.UserSession, db, tokenAuthenticator),
	}
}
