package rstokens

import (
	"context"
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/database"
)

type PostgresRSTokenStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresRSTokenStore {
	return &PostgresRSTokenStore{db: db}
}

// ----- CREATE -----

func (store *PostgresRSTokenStore) Create(ctx context.Context, token *RSToken) error {
	query := `
		INSERT INTO reset_session_tokens (token_hash, user_id, expires_at)
		VALUES ($1, $2, $3)
		RETURNING created_at
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		token.TokenHash,
		token.UserId,
		token.ExpiresAt,
	).Scan(
		&token.CreatedAt,
	)

	return database.ParseDbError(err)
}
