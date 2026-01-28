package service

import (
	"context"
	"database/sql"
)

// Transaction wrapper
func withTransaction(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	// Create transaction
	transaction, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Defer rollback (in caso di panic)
	defer func() {
		if err != nil {
			_ = transaction.Rollback() // TODO: rollback error handling?
		}
	}()

	if err = fn(transaction); err != nil {
		return err
	}

	return transaction.Commit()
}
