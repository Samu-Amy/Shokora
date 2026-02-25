package store

import (
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/store/user"
	v_token "github.com/Samu-Amy/Shokora/internal/store/verification-token"
)

type Storage struct {
	User         user.UserRepositoryI
	Product      ProductRepositoryI
	VToken       v_token.VTokenRepositoryI
	UserSession  UserSessionI
	RefreshToken RefreshTokenRepositoryI
}

func NewPostgresStorage(db *sql.DB) *Storage {
	return &Storage{
		User:         user.NewPostgresStore(db),
		Product:      NewPostgresProductStore(db),
		VToken:       v_token.NewPostgresStore(db),
		UserSession:  NewPostgresUserSessionStore(db),
		RefreshToken: NewPostgresRefreshTokenStore(db),
	}
}
