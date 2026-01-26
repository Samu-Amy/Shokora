package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
)

// ----- CREATE USER -----

func (store *PostgresUserStore) CreateUserAndSendVerification(ctx context.Context, user *User, verificationTokens *auth.VerificationTokens) error {
	// Transaction wrapper
	return withTransaction(store.db, ctx, func(transaction *sql.Tx) error {
		// Create user
		if err := store.Create(ctx, transaction, user); err != nil {
			return err
		}

		// Create verification
		if err := store.createEmailVerification(ctx, transaction, hashedToken, verificationExp, user.Id); err != nil {
			return err
		}

		return nil
	})
}

// ----- VERIFY EMAIL  -----

func (store *PostgresUserStore) VerifyEmail(ctx context.Context, plainToken string) error { // TODO: passare plain token e verificare con funzione util (?)
	return withTransaction(store.db, ctx, func(transaction *sql.Tx) error {
		// Find user related to the token
		user, err := store.getUserFromEmailVerificationToken(ctx, transaction, plainToken)
		if err != nil {
			return err
		}

		// Update user (email verified)
		user.IsVerified = true
		if err := store.setUserIsVerified(ctx, transaction, user.Id); err != nil {
			return err
		}

		// Clean email verification token
		if err := store.deleteEmailVerificationToken(ctx, transaction, user.Id); err != nil {
			return err
		}

		return nil
	})
}

func (store *PostgresUserStore) ResendEmailVerificationEmail(ctx context.Context, email string) error {
	// TODO: implementa re-invio email con token
	return nil
}

// ----- DELETE -----

func (store *PostgresUserStore) DeleteUserAndEmailVerificationToken(ctx context.Context, userId int64) error {
	return withTransaction(store.db, ctx, func(transaction *sql.Tx) error {
		if err := store.Delete(ctx, transaction, userId); err != nil {
			return err
		}

		if err := store.deleteEmailVerificationToken(ctx, transaction, userId); err != nil {
			return err
		}

		return nil
	})
}

// ----- PRIVATES -----

func (store *PostgresUserStore) createEmailVerification(ctx context.Context, transaction *sql.Tx, verificationTokens *auth.VerificationTokens, userId int64) error {
	query := `
		INSERT INTO verification_tokens (token, user_id, expiry)
		VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	_, err := transaction.ExecContext(
		ctx,
		query,
		hashedToken,
		userId,
		time.Now().Add(verificationExp),
	)

	if err != nil {
		return err
	}

	return nil
}

func (store *PostgresUserStore) getUserFromEmailVerificationToken(ctx context.Context, transaction *sql.Tx, plainToken string) (*User, error) {
	query := `
	SELECT u.id, u.first_name, u.last_name, u.email, u.is_verified, u.created_at, u.updated_at
	FROM users u
	JOIN email_verification_tokens e ON u.id = e.user_id
	WHERE e.token = $1 AND e.expiry > $2
	`
	// TODO: nel caso sia scaduto (bisogna fare un controllo separato ed eliminare sia token che user - mandare un errore ErrExpired)?

	ctx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	user := &User{}

	err := transaction.QueryRowContext(
		ctx,
		query,
		HashToken(plainToken),
		time.Now(),
	).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.IsVerified,
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

	return user, nil
}

func (store *PostgresUserStore) setUserIsVerified(ctx context.Context, transaction *sql.Tx, userId int64) error {
	query := `
		UPDATE users
		SET is_verified = true
		WHERE id = $1
	`

	_, err := transaction.ExecContext(ctx, query, userId)
	if err != nil {
		// TODO: migliorare error handling (?)
		return err
	}

	return nil
}

func (store *PostgresUserStore) deleteEmailVerificationToken(ctx context.Context, transaction *sql.Tx, userId int64) error {
	query := `
		DELETE FROM email_verification_tokens
		WHERE user_id = $1
	`

	_, err := transaction.ExecContext(ctx, query, userId)
	if err != nil {
		// TODO: migliorare error handling (?)
		return err
	}

	return nil
}
