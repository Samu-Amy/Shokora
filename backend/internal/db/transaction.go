package db

import (
	"context"
	"database/sql"
)

// - Functions -

// Transaction wrapper
func WithTransaction(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) (err error) {
	// Create transaction
	transaction, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Defer rollback (in caso di panic)
	defer func() {
		if p := recover(); p != nil {
			_ = transaction.Rollback()
			panic(p)
		} else if err != nil {
			_ = transaction.Rollback()
		}
	}()

	if err = fn(transaction); err != nil {
		return err
	}

	return transaction.Commit()
}
