package service

import (
	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/config"
	"github.com/Samu-Amy/Shokora/internal/database"
	"github.com/Samu-Amy/Shokora/internal/mailer"
	authservice "github.com/Samu-Amy/Shokora/internal/service/auth"
	userservice "github.com/Samu-Amy/Shokora/internal/service/user"
	"github.com/Samu-Amy/Shokora/internal/store"
	"go.uber.org/zap"
)

/*
The Service layer, it interacts with "internal services" (such as mailer and authenticators) and repository layer (store).
It manages the logic and internal errors, returning data and domain errors to the handlers.
*/
type Service struct {
	Auth *authservice.AuthService // Verification Tokens, Refresh Tokens, Permissions
	User *userservice.UserService // (settings, stats, achievements)
	// Menu (menu sections -> productsId)
	// Shop
	// Orders
}

func NewService(txManager database.ITransactionManager, store *store.Storage, mailer mailer.IClient, logger *zap.SugaredLogger, jwtAuthenticator *auth.JWTAuthenticator, tokenAuthenticator *auth.TokenAuthenticator, authServiceConfig config.AuthServiceConfig) *Service {
	return &Service{
		Auth: authservice.NewService(txManager, store.User, store.VToken, store.RefreshToken, store.UserSession, mailer, logger, jwtAuthenticator, tokenAuthenticator, authServiceConfig),
		User: userservice.NewService(txManager, store.User, logger),
	}
}
