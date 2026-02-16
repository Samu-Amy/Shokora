package store

import (
	"context"
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/auth"
)

type PostgresRefreshTokensStore struct {
	db *sql.DB
}

func NewPostgresRefreshTokenStore(db *sql.DB) *PostgresRefreshTokensStore {
	return &PostgresRefreshTokensStore{db: db}
}

// ----- CREATE -----

func (store *PostgresRefreshTokensStore) CreateToken(ctx context.Context, refreshToken auth.RefreshToken) error {
	query := `

	`

	queryCtx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	_, err := store.db.ExecContext(
		queryCtx,
		query,
		// TODO: continua (usa now().Add(Exp) per expires_at)
	)

	if err != nil {
		return err
	}

	return nil
}
