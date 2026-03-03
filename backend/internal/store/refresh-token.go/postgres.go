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

func (store *PostgresRefreshTokenStore) Create(ctx context.Context, transaction *sql.Tx, refreshToken *RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (session_id, token_hash, expires_at, replaces)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	err := transaction.QueryRowContext(
		queryCtx,
		query,
		refreshToken.SessionId,
		refreshToken.TokenHash,
		refreshToken.ExpiresAt,
		refreshToken.Replaces,
	).Scan(
		&refreshToken.CreatedAt,
	)

	return database.ParseDbError(err)
}

// ----- GET -----

func (store *PostgresRefreshTokenStore) GetByToken(ctx context.Context, transaction *sql.Tx, hashedToken []byte) (*TokenAndSessionData, error) {
	query := `
		SELECT r.id, r.session_id, r.expires_at, r.revoked_at, s.user_id, s.expires_at
		FROM refresh_tokens r
		JOIN user_sessions s ON r.session_id = s.id
		WHERE token_hash = $1
		FOR UPDATE;
	` //? FOR UPDATE blocca la riga fino a fine transaction (commit o rollback) - solitamente usato per get e poi update

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	var tokenAndSessionData TokenAndSessionData

	err := transaction.QueryRowContext(
		queryCtx,
		query,
		hashedToken,
	).Scan(
		&tokenAndSessionData.Id,
		&tokenAndSessionData.SessionId,
		&tokenAndSessionData.ExpiresAt,
		&tokenAndSessionData.RevokedAt,
		&tokenAndSessionData.UserId,
		&tokenAndSessionData.SessionExpiresAt,
	)

	return &tokenAndSessionData, database.ParseDbError(err)
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
