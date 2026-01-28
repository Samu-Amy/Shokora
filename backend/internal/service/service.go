package service

import (
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/store"
)

type Service struct {
	Auth *AuthService // Verification Tokens, Refresh Tokens, Permissions
	// Users (settings, stats, achievements)
	// Menu (menu sections -> productsId)
	// Shop
	// Orders
}

func NewService(db *sql.DB, store *store.Storage, tokenAuthenticator *auth.TokenAuthenticator) *Service {
	return &Service{
		Auth: NewAuthService(store.User, store.VTokens, db, tokenAuthenticator),
	}
}
