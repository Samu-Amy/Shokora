package store

import (
	"context"
	"database/sql"
)

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

// TODO: aggiungi version a user (for update check)

// ----- CREATE -----

func (store *PostgresUserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (first_name, last_name, email, password, image_url, birth_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	queryCtx, cancel := context.WithTimeout(ctx, MEDIUM_QUERY_TIMEOUT)
	defer cancel()

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.ImageUrl,
		user.BirthDate,
	).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return parseDbError(err)
}

// ----- GET -----

func (store *PostgresUserStore) GetById(ctx context.Context, userId int64) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, is_verified, user_role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	queryCtx, cancel := context.WithTimeout(ctx, MEDIUM_QUERY_TIMEOUT)
	defer cancel()

	var user User

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		userId,
	).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.IsVerified,
		&user.Role,
		// &user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return &user, parseDbError(err)
}

func (store *PostgresUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, is_verified, user_role, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_verified = true
	` // TODO: gestione verified (per chi non lo è ma accede per farsi re-inviare la mail o eliminare l'account)

	queryCtx, cancel := context.WithTimeout(ctx, MEDIUM_QUERY_TIMEOUT)
	defer cancel()

	var user User

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		email,
	).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.IsVerified,
		&user.Role,
		// &user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	// if err != nil {
	// 	switch {
	// 	case errors.Is(err, sql.ErrNoRows):
	// 		return nil, errorcodes.ErrUnauthorized // TODO: gestisci (nel caso esista ma non verificato non dovrebbe essere "not found" -> controllare dopo (nel service) is_verified (?))
	// 	default:
	// 		return nil, err
	// 	}
	// }

	return &user, parseDbError(err)
}

// ----- UPDATE -----

func (store *PostgresUserStore) Verify(ctx context.Context, userId int64) error {
	query := `
		UPDATE users
		SET is_verified = true
		WHERE id = $1
	`

	queryCtx, cancel := context.WithTimeout(ctx, MEDIUM_QUERY_TIMEOUT)
	defer cancel()

	_, err := store.db.ExecContext(queryCtx, query, userId)

	return parseDbError(err)
}

// ----- DELETE -----

func (store *PostgresUserStore) Delete(ctx context.Context, transaction *sql.Tx, userId int64) error {
	query := `DELETE FROM users WHERE id = $1`

	queryCtx, cancel := context.WithTimeout(ctx, MEDIUM_QUERY_TIMEOUT)
	defer cancel()

	_, err := transaction.ExecContext(queryCtx, query, userId)

	// TODO: gestire id non trovato (?)

	return parseDbError(err)
}
