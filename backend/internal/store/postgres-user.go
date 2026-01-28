package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
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
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	err := store.db.QueryRowContext(
		ctx,
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

	if err != nil {
		// TODO: sistema (?)
		switch {
		case strings.Contains(err.Error(), "email_key"):
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

// ----- GET -----

func (store *PostgresUserStore) GetById(ctx context.Context, userId int64) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, is_verified, user_role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	var user User

	err := store.db.QueryRowContext(
		ctx,
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

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (store *PostgresUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, password, is_verified, user_role, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_verified = true
	` // TODO: gestione verified (per chi non lo è ma accede per farsi re-inviare la mail o eliminare l'account)

	ctx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	var user User

	err := store.db.QueryRowContext(
		ctx,
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

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrUnauthorized // TODO: gestisci (nel caso esista ma non verificato non dovrebbe essere "not found" -> controllare dopo is_verified (?))
		default:
			return nil, err
		}
	}

	return &user, nil
}

// ----- DELETE -----

// TODO: come gestire l'id che rimane referenziato?
func (store *PostgresUserStore) Delete(ctx context.Context, transaction *sql.Tx, userId int64) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	_, err := transaction.ExecContext(ctx, query, userId)
	if err != nil {
		return err
	}

	return nil
}
