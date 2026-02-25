package authservice

import (
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/store"
	r_token "github.com/Samu-Amy/Shokora/internal/store/refresh-token.go"
	"github.com/Samu-Amy/Shokora/internal/store/user"
	v_token "github.com/Samu-Amy/Shokora/internal/store/verification-token"
)

type AuthService struct {
	userRepo           user.UserRepositoryI // TODO: serve (tolta creazione utente)?
	vTokenRepo         v_token.VTokenRepositoryI
	refreshTokenRepo   r_token.RefreshTokenRepositoryI
	userSessionRepo    store.UserSessionI
	db                 *sql.DB
	tokenAuthenticator *auth.TokenAuthenticator
}

func NewService(userRepo user.UserRepositoryI, vTokensRepo v_token.VTokenRepositoryI, refreshTokensRepo r_token.RefreshTokenRepositoryI, userSessionRepo store.UserSessionI, db *sql.DB, tokenAuthenticator *auth.TokenAuthenticator) *AuthService {
	return &AuthService{userRepo, vTokensRepo, refreshTokensRepo, userSessionRepo, db, tokenAuthenticator}
}
