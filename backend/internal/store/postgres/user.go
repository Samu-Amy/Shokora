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

func (store *PostgresUserStore) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (first_name, last_name, email, password)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	err := store.db.QueryRowContext(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

// TODO: aggiungi eliminazione account (come gestire l'id che rimane referenziato?)
