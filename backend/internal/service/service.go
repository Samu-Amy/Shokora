package service

import (
	"context"
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/store"
)

type Service struct {
	Auth *AuthService // Verification Tokens, Refresh Tokens, Permissions
	// Users (settings, stats, achievements)
	// Menu (menu sections -> productsId)
	// Shop
	// Orders
}

func NewService(db *sql.DB, store *store.Storage) Service {
	return Service{
		Auth: NewAuthService(store.User, store.VTokens, db),
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
	// if err := fn(transaction); err != nil {
	// 	_ = transaction.Rollback() // TODO: rollback error handling?
	// 	return err
	// }

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
