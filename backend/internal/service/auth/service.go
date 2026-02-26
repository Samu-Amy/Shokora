package authservice

import (
	"github.com/Samu-Amy/Shokora/internal/api"
	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/database"
	"github.com/Samu-Amy/Shokora/internal/mailer"
	rtoken "github.com/Samu-Amy/Shokora/internal/store/refresh-token.go"
	"github.com/Samu-Amy/Shokora/internal/store/user"
	session "github.com/Samu-Amy/Shokora/internal/store/user-session"
	v_token "github.com/Samu-Amy/Shokora/internal/store/verification-token"
	"go.uber.org/zap"
)

type AuthService struct {
	txManager          database.ITransactionManager
	userRepo           user.IUserRepository
	vTokenRepo         v_token.IVTokenRepository
	refreshTokenRepo   rtoken.IRefreshTokenRepository
	userSessionRepo    session.IUserSessionRepository
	mailer             mailer.IClient
	logger             *zap.SugaredLogger
	jwtAuthenticator   *auth.JWTAuthenticator
	tokenAuthenticator *auth.TokenAuthenticator
	config             AuthServiceConfig
}

func NewService(txManager database.ITransactionManager, userRepo user.IUserRepository, vTokenRepo v_token.IVTokenRepository, refreshTokensRepo rtoken.IRefreshTokenRepository, userSessionRepo session.IUserSessionRepository, mailer mailer.IClient, logger *zap.SugaredLogger, jwtAuthenticator *auth.JWTAuthenticator, tokenAuthenticator *auth.TokenAuthenticator, config AuthServiceConfig) *AuthService {
	return &AuthService{txManager, userRepo, vTokenRepo, refreshTokensRepo, userSessionRepo, mailer, logger, jwtAuthenticator, tokenAuthenticator, config}
}

// - Config -
type AuthServiceConfig struct {
	PasswordHashingCost int
	Token               api.TokenConfig
	Mail                MailConfig
}

type MailConfig struct {
	IsSandboxEnv bool
	FrontEndURL  string
}
