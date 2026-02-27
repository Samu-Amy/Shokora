package session

import (
	"context"
	"database/sql"
	"time"

	"github.com/Samu-Amy/Shokora/internal/database"
)

type PostgresUserSessionStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresUserSessionStore {
	return &PostgresUserSessionStore{db: db}
}

// ----- CREATE -----

func (store *PostgresUserSessionStore) Create(ctx context.Context, transaction *sql.Tx, session *UserSession, sessionExp time.Duration) error {
	query := `
		INSERT INTO user_sessions (user_id, expires_at)
		VALUES ($1, NOW() + $2)
		RETURNING id, expires_at, created_at
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	err := transaction.QueryRowContext(
		queryCtx,
		query,
		session.UserId,
		sessionExp,
	).Scan(
		&session.Id,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	return err
}

// ----- DELETE -----

func (store *PostgresUserSessionStore) Delete(ctx context.Context, transaction *sql.Tx, sessionId int64) error {
	query := `DELETE FROM user_sessions WHERE id = $1`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	return database.HandleExecContextResult(transaction.ExecContext(queryCtx, query, sessionId))
}
