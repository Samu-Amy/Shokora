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

func (store *PostgresUserSessionStore) Create(ctx context.Context, transaction *sql.Tx, userId int64, sessionExp time.Duration) (int64, error) {
	query := `
		INSERT INTO user_sessions (user_id, expires_at)
		VALUES ($1, $2)
		RETURNING id
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	var sessionId int64

	err := transaction.QueryRowContext(
		queryCtx,
		query,
		userId,
		time.Now().Add(sessionExp).UTC(),
	).Scan(
		&sessionId,
	)

	return sessionId, database.ParseDbError(err)
}

// ----- DELETE -----

// - Delete a session by its id -

func (store *PostgresUserSessionStore) Delete(ctx context.Context, sessionId int64) error {
	query := `DELETE FROM user_sessions WHERE id = $1`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	return database.HandleExecContextResult(store.db.ExecContext(queryCtx, query, sessionId))
}

// - Delete expired session -

func (store *PostgresUserSessionStore) DeleteExpired(ctx context.Context, userId int64) error {
	query := `
		DELETE FROM user_sessions
		WHERE user_id = $1 AND expires_at < NOW()
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	return database.HandleExecContextResult(store.db.ExecContext(queryCtx, query, userId))
}

// - Delete all sessions for one user by the userId -

func (store *PostgresUserSessionStore) DeleteOtherUserSessions(ctx context.Context, transaction *sql.Tx, userId, sessionId int64) error {
	query := `
		DELETE FROM user_sessions
		WHERE user_id = $1 AND id != $2
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	return database.HandleExecContextResult(transaction.ExecContext(queryCtx, query, userId, sessionId))
}
