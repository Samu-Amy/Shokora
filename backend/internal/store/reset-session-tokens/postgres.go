package rstoken

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

// ----- GET -----

func (store *PostgresRSTokenStore) Get(ctx context.Context, hashedToken []byte) (*RSToken, error) {
	query := `
		SELECT user_id, expires_at, created_at
		FROM reset_session_tokens
		WHERE token_hash = $1
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	var rsToken RSToken

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		hashedToken,
	).Scan(
		&rsToken.UserId,
		&rsToken.ExpiresAt,
		&rsToken.CreatedAt,
	)

	if err != nil {
		return nil, database.ParseDbError(err)
	}

	return &rsToken, nil
}

// ----- DELETE -----

func (store *PostgresRSTokenStore) Delete(ctx context.Context, hashedToken []byte) error {
	query := `DELETE FROM user_sessions WHERE token_hash = $1`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	return database.HandleExecContextResult(store.db.ExecContext(queryCtx, query, hashedToken))
}
