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
	userRepo           user.UserRepositoryI // TODO: serve (tolta creazione utente)?
	vTokenRepo         v_token.VTokenRepositoryI
	refreshTokenRepo   rtoken.RefreshTokenRepositoryI
	userSessionRepo    session.UserSessionI
	db                 *sql.DB
	tokenAuthenticator *auth.TokenAuthenticator
}

func NewService(userRepo user.UserRepositoryI, vTokensRepo v_token.VTokenRepositoryI, refreshTokensRepo rtoken.RefreshTokenRepositoryI, userSessionRepo session.UserSessionI, db *sql.DB, tokenAuthenticator *auth.TokenAuthenticator) *AuthService {
	return &AuthService{userRepo, vTokensRepo, refreshTokensRepo, userSessionRepo, db, tokenAuthenticator}
}
