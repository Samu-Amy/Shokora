package store

import (
	"context"
	"database/sql"
	"time"
)

// - CREATE -

func (store *PostgresUserStore) CreateAndSendVerification(ctx context.Context, user *User, token string, verificationExp time.Duration) error {
	// Transaction wrapper
	return withTransaction(store.db, ctx, func(transaction *sql.Tx) error {
		// Create user
		if err := store.Create(ctx, transaction, user); err != nil {
			return err
		}

		// Create verification
		if err := store.createUserVerification(ctx, transaction, token, verificationExp, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (store *PostgresUserStore) createUserVerification(ctx context.Context, transaction *sql.Tx, token string, verificationExp time.Duration, userId int64) error {
	query := `
		INSERT INTO user_verification_tokens (token, user_id, expiry)
		VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	_, err := transaction.ExecContext(
		ctx,
		query,
		token,
		userId,
		time.Now().Add(verificationExp),
	)

	if err != nil {
		return err
	}

	return nil
}
