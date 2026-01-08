package postgres

import (
	"context"
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/store/models"
)

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

func (s *PostgresUserStore) Create(context.Context, *models.User) error {
	return nil
}
