package r_token

import (
	"context"
	"database/sql"
	"time"

	"github.com/Samu-Amy/Shokora/internal/db"
)

type PostgresRefreshTokenStore struct {
	db *sql.DB
}

func NewPostgresRefreshTokenStore(db *sql.DB) *PostgresRefreshTokenStore {
	return &PostgresRefreshTokenStore{db: db}
}

// ----- CREATE -----

func (store *PostgresRefreshTokenStore) CreateToken(ctx context.Context, transaction *sql.Tx, refreshToken *RefreshToken, tokenExp time.Duration) error {
	query := `
		INSERT INTO refresh_tokens (session_id, token_hash, expires_at, session_exp, replaces)
		VALUES ($1, $2, NOW() + $3, NOW() + $4, $5)
		RETURNING expires_at, created_at
	`

	queryCtx, cancel := context.WithTimeout(ctx, db.MEDIUM_QUERY_TIMEOUT)
	defer cancel()

	err := transaction.QueryRowContext(
		queryCtx,
		query,
		refreshToken.SessionId,
		refreshToken.TokenHash,
		tokenExp, // TODO: questo deve essere di quanto aumenta la sessione (es. 7 giorni) e si aggiunge a NOW (anche se forse così non la si sta aumentando molto, se scade tra 7 giorni ed aggiungo 7 giorni da ora, scade comunque tra 7 giorni) (se non supera la durata massima della sessione)
		refreshToken.Replaces,
	).Scan(
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
	)

	return db.ParseDbError(err)
}

// ----- GET -----

func (store *PostgresRefreshTokenStore) GetToken(ctx context.Context, transaction *sql.Tx, hashedToken []byte) (*RefreshToken, error) {
	query := `
		SELECT id, user_id, session_id, token_hash, expires_at, replaces, revoked_at, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
		FOR UPDATE;
	` //? FOR UPDATE blocca la riga fino a fine transaction (commit o rollback) - solitamente usato per get e poi update

	queryCtx, cancel := context.WithTimeout(ctx, db.MEDIUM_QUERY_TIMEOUT)
	defer cancel()

	var refreshToken RefreshToken

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

func (store *PostgresRefreshTokenStore) RevokeTokenById(ctx context.Context, transaction *sql.Tx, tokenId int64, revokedAt time.Time) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = $1
		WHERE id = $2 AND revoked_at IS NULL
	`

	queryCtx, cancel := context.WithTimeout(ctx, db.MEDIUM_QUERY_TIMEOUT)
	defer cancel()

	return db.HandleExecContextResult(transaction.ExecContext(
		queryCtx,
		query,
		revokedAt,
		tokenId,
	))
}

// ----- DELETE -----

// func (store *PostgresRefreshTokenStore) DeleteSessionById(ctx context.Context, userId int64, sessionId uuid.UUID) error {
// 	query := `
// 		DELETE FROM refresh_tokens
// 		WHERE user_id = $1 AND session_id = $2
// 	`

// 	queryCtx, cancel := context.WithTimeout(ctx, db.MEDIUM_QUERY_TIMEOUT)
// 	defer cancel()

// 	return db.HandleExecContextResult(store.db.ExecContext(
// 		queryCtx,
// 		query,
// 		userId,
// 		sessionId,
// 	))
// }
