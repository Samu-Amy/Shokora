package store

import (
	"database/sql"
)

type Storage struct {
	User          UserRepositoryI
	Product       ProductRepositoryI
	VTokens       VTokensRepositoryI
	RefreshTokens RefreshTokensRepositoryI
}

func NewPostgresStorage(db *sql.DB) *Storage {
	return &Storage{
		User:          NewPostgresUserStore(db),
		Product:       NewPostgresProductStore(db),
		VTokens:       NewPostgresVTokenStore(db),
		RefreshTokens: NewPostgresRefreshTokenStore(db),
	}
}
