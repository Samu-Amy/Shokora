package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	User    UserRepository
	Product ProductRepository
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		User:    NewPostgresUserStore(db),
		Product: NewPostgresProductStore(db),
	}
}

// Transaction wrapper
func withTransaction(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	// Create transaction
	transaction, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Use transaction
	if err := fn(transaction); err != nil {
		_ = transaction.Rollback() // TODO: rollback error handling?
		return err
	}

	return transaction.Commit()
}
