package store

import (
	"database/sql"
)

type Storage struct {
	User    UserRepositoryI
	Product ProductRepositoryI
	VTokens VTokensRepositoryI
}

func NewPostgresStorage(db *sql.DB) *Storage {
	return &Storage{
		User:    NewPostgresUserStore(db),
		Product: NewPostgresProductStore(db),
		VTokens: NewPostgresVTokenStore(db),
	}
}
