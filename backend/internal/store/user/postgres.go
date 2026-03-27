package user

import (
	"context"
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/database"
)

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

// TODO: aggiungi version a user (for update check)

// ----- CREATE -----

func (store *PostgresUserStore) Create(ctx context.Context, transaction *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (first_name, last_name, email, password, birthday)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	err := transaction.QueryRowContext(
		queryCtx,
		query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		// user.ImageUrl,
		user.Birthday,
	).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return database.ParseDbError(err)
}

// ----- GET -----

func (store *PostgresUserStore) GetById(ctx context.Context, userId int64) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, birthday, is_verified, user_role, permissions, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
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
		// &user.ImageUrl,
		&user.Birthday,
		&user.IsVerified,
		&user.Role,
		&user.Permissions,
		// &user.Version,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, database.ParseDbError(err)
	}

	return &user, nil
}

func (store *PostgresUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, birthday, is_verified, user_role, permissions, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
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
		// &user.ImageUrl,
		&user.Birthday,
		&user.IsVerified,
		&user.Role,
		&user.Permissions,
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

	if err != nil {
		return nil, database.ParseDbError(err)
	}

	return &user, nil
}

func (store *PostgresUserStore) GetUserVerificationDataByEmail(ctx context.Context, transaction *sql.Tx, email string) (*UserVerificationData, error) {
	query := `
		SELECT id, first_name
		FROM users
		WHERE email = $1
		FOR UPDATE
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	var userData UserVerificationData

	err := transaction.QueryRowContext(
		queryCtx,
		query,
		email,
	).Scan(
		&userData.Id,
		&userData.FirstName,
	)

	if err != nil {
		return nil, database.ParseDbError(err)
	}

	return &userData, nil
}

func (store *PostgresUserStore) GetPassword(ctx context.Context, transaction *sql.Tx, userId int64) ([]byte, error) {
	query := `
		SELECT password
		FROM users
		WHERE id = $1
		FOR UPDATE
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	var hashedPassword []byte

	err := transaction.QueryRowContext(
		queryCtx,
		query,
		userId,
	).Scan(
		&hashedPassword,
	)

	if err != nil {
		return nil, database.ParseDbError(err)
	}

	return hashedPassword, nil
}

// ----- UPDATE -----

func (store *PostgresUserStore) UpdatePassword(ctx context.Context, transaction *sql.Tx, userId int64, hashedPassword []byte) error {
	query := `
		UPDATE users
		SET password = $1
		WHERE id = $2
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	return database.HandleExecContextResult(transaction.ExecContext(queryCtx, query, hashedPassword, userId))
}

func (store *PostgresUserStore) SetIsVerified(ctx context.Context, userId int64) error {
	query := `
		UPDATE users
		SET is_verified = true
		WHERE id = $1
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	return database.HandleExecContextResult(store.db.ExecContext(queryCtx, query, userId))
}

// ----- DELETE -----

func (store *PostgresUserStore) Delete(ctx context.Context, userId int64) error {
	query := `DELETE FROM users WHERE id = $1`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	return database.HandleExecContextResult(store.db.ExecContext(queryCtx, query, userId))
}
