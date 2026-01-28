package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// TODO: fai lookup anche per verification_type (token + verification_type | OPT + email + verification_type + attempts) -- SET otp_attempts = otp_attempts + 1 (aggiorna atomicamente attempts)

// ----- VERIFY EMAIL  -----

// func (store *PostgresUserStore) VerifyEmail(ctx context.Context, plainToken string) error { // TODO: passare plain token e verificare con funzione util (?)
// 	return withTransaction(store.db, ctx, func(transaction *sql.Tx) error {
// 		// Find user related to the token
// 		user, err := store.getUserFromEmailVerificationToken(ctx, transaction, plainToken)
// 		if err != nil {
// 			return err
// 		}

// 		// Update user (email verified)
// 		user.IsVerified = true
// 		if err := store.setUserIsVerified(ctx, transaction, user.Id); err != nil {
// 			return err
// 		}

// 		// Clean email verification token
// 		if err := store.deleteEmailVerificationToken(ctx, transaction, user.Id); err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// }

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
