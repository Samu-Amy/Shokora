package store

import (
	"database/sql"
)

type Storage struct {
	User         UserRepositoryI
	Product      ProductRepositoryI
	VToken       VTokenRepositoryI
	UserSession  UserSessionI
	RefreshToken RefreshTokenRepositoryI
}

func NewPostgresStorage(db *sql.DB) *Storage {
	return &Storage{
		User:         NewPostgresUserStore(db),
		Product:      NewPostgresProductStore(db),
		VToken:       NewPostgresVTokenStore(db),
		UserSession:  NewPostgresUserSessionStore(db),
		RefreshToken: NewPostgresRefreshTokenStore(db),
	}
}
