package authservice

import (
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/mailer"
	rtoken "github.com/Samu-Amy/Shokora/internal/store/refresh-token.go"
	"github.com/Samu-Amy/Shokora/internal/store/user"
	session "github.com/Samu-Amy/Shokora/internal/store/user-session"
	v_token "github.com/Samu-Amy/Shokora/internal/store/verification-token"
)

type AuthService struct {
	db                 *sql.DB
	mailer             mailer.IClient
	userRepo           user.IUserRepository
	vTokenRepo         v_token.IVTokenRepository
	refreshTokenRepo   rtoken.IRefreshTokenRepository
	userSessionRepo    session.IUserSessionRepository
	jwtAuthenticator   *auth.JWTAuthenticator
	tokenAuthenticator *auth.TokenAuthenticator
	config             AuthServiceConfig
}

func NewService(db *sql.DB, mailer mailer.IClient, userRepo user.IUserRepository, vTokensRepo v_token.IVTokenRepository, refreshTokensRepo rtoken.IRefreshTokenRepository, userSessionRepo session.IUserSessionRepository, jwtAuthenticator *auth.JWTAuthenticator, tokenAuthenticator *auth.TokenAuthenticator, config AuthServiceConfig) *AuthService {
	return &AuthService{db, mailer, userRepo, vTokensRepo, refreshTokensRepo, userSessionRepo, jwtAuthenticator, tokenAuthenticator, config}
}

// - Config -
type AuthServiceConfig struct {
	PasswordHashingCost int
}
