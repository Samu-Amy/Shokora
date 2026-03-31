package authservice

import (
	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/config"
	"github.com/Samu-Amy/Shokora/internal/database"
	"github.com/Samu-Amy/Shokora/internal/mailer"
	oauthstate "github.com/Samu-Amy/Shokora/internal/store/oauth-states"
	rtoken "github.com/Samu-Amy/Shokora/internal/store/refresh-token"
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
	oAuthStateRepo     oauthstate.IOAuthStateRepository
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

func NewService(txManager database.ITransactionManager, vTokenRepo v_token.IVTokenRepository, rsTokenRepo rstoken.IResetSessionTokenRepository, oAuthStateRepo oauthstate.IOAuthStateRepository, refreshTokensRepo rtoken.IRefreshTokenRepository, userSessionRepo session.IUserSessionRepository, userRepo user.IUserRepository, userSettingsRepo usersettings.IUserSettingsRepository, mailer mailer.IClient, logger *zap.SugaredLogger, jwtAuthenticator *auth.JWTAuthenticator, tokenAuthenticator *auth.TokenAuthenticator, config config.AuthServiceConfig) *AuthService {
	return &AuthService{txManager, vTokenRepo, rsTokenRepo, oAuthStateRepo, refreshTokensRepo, userSessionRepo, userRepo, userSettingsRepo, mailer, logger, jwtAuthenticator, tokenAuthenticator, config}
}
