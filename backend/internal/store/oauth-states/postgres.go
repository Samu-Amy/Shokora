package oauthstate

import (
	"context"
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/database"
)

type PostgresOAuthStateStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresOAuthStateStore {
	return &PostgresOAuthStateStore{db: db}
}

// ----- CREATE -----

func (store *PostgresOAuthStateStore) Create(ctx context.Context, state string) error {
	query := `
		INSERT INTO oauth_states (state)
		VALUES ($1)
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	return database.HandleExecContextResult(store.db.ExecContext(
		queryCtx,
		query,
		state,
	))
}

// ----- GET -----

// func (store *PostgresOAuthStateStore) Get(ctx context.Context, oAuthState *OAuthState) error {
// 	query := `
// 		SELECT created_at
// 		FROM oauth_states
// 		WHERE state = $1
// 	`

// 	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
// 	defer cancel()

// 	err := store.db.QueryRowContext(
// 		queryCtx,
// 		query,
// 		oAuthState.State,
// 	).Scan(
// 		&oAuthState.CreatedAt,
// 	)

// 	return database.ParseDbError(err)
// }

// ----- DELETE -----

func (store *PostgresOAuthStateStore) Delete(ctx context.Context, transaction *sql.Tx, state string) error {
	query := `
		DELETE FROM oauth_states
		WHERE state = $1
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	return database.HandleExecContextResult(transaction.ExecContext(
		queryCtx,
		query,
		state,
	))
}
