package store

import (
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/store/models"
	"github.com/Samu-Amy/Shokora/internal/store/postgres"
)

type Storage struct {
	User    models.UserRepository
	Product models.ProductRepository
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		User:    postgres.NewPostgresUserStore(db),
		Product: postgres.NewPostgresProductStore(db),
	}
}
