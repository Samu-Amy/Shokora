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

func (store *PostgresRefreshTokensStore) CreateToken(ctx context.Context, queryer Queryer, refreshToken *auth.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, session_id, token_hash, expires_at, replaces)
		VALUES ($1, $2, $3, NOW() + $4, $5)
		RETURNING expires_at, created_at
	`

	queryCtx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	err := queryer.QueryRowContext(
		queryCtx,
		query,
		refreshToken.UserId,
		refreshToken.SessionId,
		refreshToken.HashedToken,
		refreshToken.Exp,
		refreshToken.Replaces,
	).Scan(
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

// ----- GET -----

func (store *PostgresRefreshTokensStore) GetToken(ctx context.Context, transaction *sql.Tx, hashedToken []byte) (*RefreshTokens, error) {
	query := `
		SELECT id, user_id, session_id, token_hash, expires_at, replaces, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
		FOR UPDATE;
	` //? FOR UPDATE blocca la riga fino a fine transaction (commit o rollback) - solitamente usato per get e poi update

	queryCtx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	var refreshToken RefreshTokens

	err := transaction.QueryRowContext(
		queryCtx,
		query,
		hashedToken,
	).Scan(
		&refreshToken.Id,
		&refreshToken.UserId,
		&refreshToken.SessionId,
		&refreshToken.TokenHash,
		&refreshToken.Exp,
		&refreshToken.Replaces,
		&refreshToken.RevokedAt,
		&refreshToken.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

// ----- UPDATE -----

func (store *PostgresRefreshTokensStore) RevokeTokenById(ctx context.Context, transaction *sql.Tx, tokenId int64, revokedAt time.Time) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = $1
		WHERE id = $2 AND revoked_at = NULL
	`

	queryCtx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	_, err := transaction.ExecContext(
		queryCtx,
		query,
		revokedAt,
		tokenId,
	)

	if err != nil {
		return err
	}

	return nil
}
