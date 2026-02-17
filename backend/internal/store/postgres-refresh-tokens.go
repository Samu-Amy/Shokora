package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
)

type PostgresRefreshTokensStore struct {
	db *sql.DB
}

func NewPostgresRefreshTokenStore(db *sql.DB) *PostgresRefreshTokensStore {
	return &PostgresRefreshTokensStore{db: db}
}

// ----- CREATE -----

func (store *PostgresRefreshTokensStore) CreateToken(ctx context.Context, refreshToken auth.RefreshToken) (*time.Time, error) {
	query := `
		INSERT INTO refresh_tokens (user_id, session_id, token_hash, expires_at, replaces)
		VALUES ($1, $2, $3, NOW() + $4, $5)
		RETURNING expires_at
	`

	queryCtx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	var tokenExpiresAt time.Time

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		refreshToken.UserId,
		refreshToken.SessionId,
		refreshToken.HashedToken,
		refreshToken.Exp,
		refreshToken.Replaces, // If replaces != nil -> rotation, else is a new token
	).Scan(
		&tokenExpiresAt,
	)

	if err != nil {
		return nil, err
	}

	return &tokenExpiresAt, nil
}
