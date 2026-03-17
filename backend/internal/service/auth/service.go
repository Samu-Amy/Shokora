package authservice

import (
	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/config"
	"github.com/Samu-Amy/Shokora/internal/database"
	"github.com/Samu-Amy/Shokora/internal/mailer"
	rtoken "github.com/Samu-Amy/Shokora/internal/store/refresh-token.go"
	rstoken "github.com/Samu-Amy/Shokora/internal/store/reset-session-tokens"
	"github.com/Samu-Amy/Shokora/internal/store/user"
	session "github.com/Samu-Amy/Shokora/internal/store/user-session"
	usersettings "github.com/Samu-Amy/Shokora/internal/store/user-settings"
	v_token "github.com/Samu-Amy/Shokora/internal/store/verification-token"
	"go.uber.org/zap"
)

type AuthService struct {
	txManager          database.ITransactionManager
	vTokenRepo         v_token.IVTokenRepository
	rsTokenRepo        rstoken.IResetSessionTokenRepository
	refreshTokenRepo   rtoken.IRefreshTokenRepository
	userSessionRepo    session.IUserSessionRepository
	userRepo           user.IUserRepository
	userSettingsRepo   usersettings.IUserSettingsRepository
	mailer             mailer.IClient
	logger             *zap.SugaredLogger
	jwtAuthenticator   *auth.JWTAuthenticator
	tokenAuthenticator *auth.TokenAuthenticator
	config             config.AuthServiceConfig
}

func NewService(txManager database.ITransactionManager, vTokenRepo v_token.IVTokenRepository, rsTokenRepo rstoken.IResetSessionTokenRepository, refreshTokensRepo rtoken.IRefreshTokenRepository, userSessionRepo session.IUserSessionRepository, userRepo user.IUserRepository, userSettingsRepo usersettings.IUserSettingsRepository, mailer mailer.IClient, logger *zap.SugaredLogger, jwtAuthenticator *auth.JWTAuthenticator, tokenAuthenticator *auth.TokenAuthenticator, config config.AuthServiceConfig) *AuthService {
	return &AuthService{txManager, vTokenRepo, rsTokenRepo, refreshTokensRepo, userSessionRepo, userRepo, userSettingsRepo, mailer, logger, jwtAuthenticator, tokenAuthenticator, config}
}
