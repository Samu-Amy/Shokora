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
		INSERT INTO users (google_id, first_name, last_name, email, password, birthday, is_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	isVerified := false
	if user.IsVerified {
		if user.GoogleId != "" {
			// OAuth Google -> verified
			isVerified = true
		} else {
			// Fix possible error (user can't be verified at creation without OAuth)
			user.IsVerified = false
		}
	}

	err := transaction.QueryRowContext(
		queryCtx,
		query,
		user.GoogleId,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		// user.ImageUrl,
		user.Birthday,
		isVerified,
	).Scan(
		&user.Id,
		&user.Role,
		&user.Permissions,
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

func (store *PostgresUserStore) GetByGoogleId(ctx context.Context, googleId string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, birthday, is_verified, user_role, permissions, created_at, updated_at
		FROM users
		WHERE google_id = $1
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	var user User

	err := store.db.QueryRowContext(
		queryCtx,
		query,
		googleId,
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

func (store *PostgresUserStore) GetByEmailForUpdate(ctx context.Context, transaction *sql.Tx, email string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, birthday, is_verified, user_role, permissions, created_at, updated_at
		FROM users
		WHERE email = $1
		FOR UPDATE
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	var user User

	err := transaction.QueryRowContext(
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

	if err != nil {
		return nil, database.ParseDbError(err)
	}

	return &user, nil
}

func (store *PostgresUserStore) GetUserVerificationDataByEmailForUpdate(ctx context.Context, transaction *sql.Tx, email string) (*UserVerificationData, error) {
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

func (store *PostgresUserStore) GetPasswordForUpdate(ctx context.Context, transaction *sql.Tx, userId int64) ([]byte, error) {
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

func (store *PostgresUserStore) SetGoogleId(ctx context.Context, transaction *sql.Tx, userId int64, googleId string) error {
	query := `
		UPDATE users
		SET google_id = $1
		WHERE id = $2
	`

	queryCtx, cancel := context.WithTimeout(ctx, database.MediumQueryTimeout)
	defer cancel()

	return database.HandleExecContextResult(transaction.ExecContext(queryCtx, query, googleId, userId))
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
