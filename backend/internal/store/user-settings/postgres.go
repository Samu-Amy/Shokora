package usersettings

import (
	"context"
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/database"
)

type PostgresUserSettingsStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresUserSettingsStore {
	return &PostgresUserSettingsStore{db: db}
}

// ----- CREATE -----

func (store *PostgresUserSettingsStore) Create(ctx context.Context, transaction *sql.Tx, userId int64) (int64, error) {
	query := `
		INSERT INTO user_settings (user_id)
		VALUES ($1)
		RETURNING id
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	var settingsId int64

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		userId,
	).Scan(
		&settingsId,
	)

	return settingsId, database.ParseDbError(err)
}

// ----- UPDATE -----

func (store *PostgresUserSettingsStore) Update(ctx context.Context, settings *UserSettings) error {
	query := `
		UPDATE user_settings
		SET two_factor_auth = $1
		WHERE id = $2
		RETURNING updated_at
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		settings.HasTwoFactorAuth,
	).Scan(
		&settings.UpdatedAt,
	)

	return database.ParseDbError(err)
}

// ----- DELETE -----

func (store *PostgresUserSettingsStore) Delete(ctx context.Context, settingsId int64) error {
	query := `DELETE FROM user_settings WHERE id = $1`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	return database.HandleExecContextResult(store.db.ExecContext(queryCtx, query, settingsId))
}
