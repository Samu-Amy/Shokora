package store

import (
	"database/sql"
)

type Storage struct {
	User    UserRepository
	Product ProductRepository
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		User:    NewPostgresUserStore(db),
		Product: NewPostgresProductStore(db),
	}
}
