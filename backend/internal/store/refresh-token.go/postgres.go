package rtoken

import (
	"context"
	"database/sql"
	"time"

	"github.com/Samu-Amy/Shokora/internal/database"
)

type PostgresRefreshTokenStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresRefreshTokenStore {
	return &PostgresRefreshTokenStore{db: db}
}

// ----- CREATE -----

func (store *PostgresRefreshTokenStore) Create(ctx context.Context, transaction *sql.Tx, refreshToken *RefreshToken, tokenExp time.Duration) error {
	query := `
		INSERT INTO refresh_tokens (session_id, token_hash, expires_at, replaces)
		VALUES ($1, $2, $3, $4)
		RETURNING expires_at
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	err := transaction.QueryRowContext(
		queryCtx,
		query,
		refreshToken.SessionId,
		refreshToken.TokenHash,
		time.Now().Add(tokenExp),
		refreshToken.Replaces,
	).Scan(
		&refreshToken.ExpiresAt,
	)

	return database.ParseDbError(err)
}

// ----- GET -----

func (store *PostgresRefreshTokenStore) Get(ctx context.Context, transaction *sql.Tx, hashedToken []byte) (*RefreshToken, error) {
	query := `
		SELECT id, session_id, token_hash, expires_at, replaces, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
		FOR UPDATE;
	` //? FOR UPDATE blocca la riga fino a fine transaction (commit o rollback) - solitamente usato per get e poi update

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	var refreshToken RefreshToken

	err := transaction.QueryRowContext(
		queryCtx,
		query,
		hashedToken,
	).Scan(
		&refreshToken.Id,
		&refreshToken.SessionId,
		&refreshToken.TokenHash,
		&refreshToken.ExpiresAt,
		&refreshToken.Replaces,
		&refreshToken.RevokedAt,
		&refreshToken.CreatedAt,
	)

	return &refreshToken, database.ParseDbError(err)
}

// ----- UPDATE -----

func (store *PostgresRefreshTokenStore) RevokeById(ctx context.Context, transaction *sql.Tx, tokenId int64, revokedAt time.Time) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = $1
		WHERE id = $2 AND revoked_at IS NULL
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	return database.HandleExecContextResult(transaction.ExecContext(
		queryCtx,
		query,
		revokedAt,
		tokenId,
	))
}
