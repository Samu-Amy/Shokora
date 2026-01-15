package store

import (
	"context"
	"database/sql"
	"errors"
)

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

// TODO: add version (for update check)

// ----- CREATE -----

func (store *PostgresUserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (first_name, last_name, email, password)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

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

// ----- GET -----

func (store *PostgresUserStore) GetById(ctx context.Context, userId int64) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	var user User

	err := store.db.QueryRowContext(
		ctx,
		query,
		userId,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		// &user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

// TODO: aggiungi eliminazione account (come gestire l'id che rimane referenziato?)
