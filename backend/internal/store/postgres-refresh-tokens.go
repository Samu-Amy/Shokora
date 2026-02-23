package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/errorcodes"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
		// Reuse detection
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" && pgErr.Constraint == "refresh_tokens_replaces_unique" {
				return errorcodes.InternalErrReusedToken
			}
		}
		return err
	}

	return nil
}

// ----- GET -----

func (store *PostgresRefreshTokensStore) GetToken(ctx context.Context, transaction *sql.Tx, hashedToken []byte) (*auth.RefreshToken, error) {
	query := `
		SELECT id, user_id, session_id, token_hash, expires_at, replaces, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
		FOR UPDATE;
	` //? FOR UPDATE blocca la riga fino a fine transaction (commit o rollback) - solitamente usato per get e poi update

	queryCtx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	var refreshToken auth.RefreshToken

	err := transaction.QueryRowContext(
		queryCtx,
		query,
		hashedToken,
	).Scan(
		&refreshToken.Id,
		&refreshToken.UserId,
		&refreshToken.SessionId,
		&refreshToken.HashedToken,
		&refreshToken.ExpiresAt,
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
		WHERE id = $2 AND revoked_at IS NULL
	`

	queryCtx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	res, err := transaction.ExecContext(
		queryCtx,
		query,
		revokedAt,
		tokenId,
	)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errorcodes.InternalErrTokenNotFoundOrAlreadyRevoked // Token already revoked or not found
	}

	return nil
}

// ----- DELETE -----

func (store *PostgresRefreshTokensStore) DeleteSessionById(ctx context.Context, userId int64, sessionId uuid.UUID) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE user_id = $1 AND session_id = $2
	`

	queryCtx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	_, err := store.db.ExecContext(
		queryCtx,
		query,
		userId,
		sessionId,
	)

	if err != nil {
		return err
	}

	return nil
}
