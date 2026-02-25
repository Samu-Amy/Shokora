package authservice

import (
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/auth"
	rtoken "github.com/Samu-Amy/Shokora/internal/store/refresh-token.go"
	"github.com/Samu-Amy/Shokora/internal/store/user"
	session "github.com/Samu-Amy/Shokora/internal/store/user-session"
	v_token "github.com/Samu-Amy/Shokora/internal/store/verification-token"
)

type AuthService struct {
	userRepo           user.IUserRepository
	vTokenRepo         v_token.IVTokenRepository
	refreshTokenRepo   rtoken.IRefreshTokenRepository
	userSessionRepo    session.IUserSessionRepository
	db                 *sql.DB
	tokenAuthenticator *auth.TokenAuthenticator
}

func NewService(userRepo user.IUserRepository, vTokensRepo v_token.IVTokenRepository, refreshTokensRepo rtoken.IRefreshTokenRepository, userSessionRepo session.IUserSessionRepository, db *sql.DB, tokenAuthenticator *auth.TokenAuthenticator) *AuthService {
	return &AuthService{userRepo, vTokensRepo, refreshTokensRepo, userSessionRepo, db, tokenAuthenticator}
}
