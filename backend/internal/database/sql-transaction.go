package database

import (
	"context"
	"database/sql"
)

type SQLTransactionManager struct {
	db *sql.DB
}

func NewSQLTransactionManager(db *sql.DB) *SQLTransactionManager {
	return &SQLTransactionManager{db: db}
}

func (txManager *SQLTransactionManager) WithTx(ctx context.Context, fn func(*sql.Tx) error) (err error) {
	// Create transaction
	transaction, err := txManager.db.BeginTx(ctx, nil)
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

	// Use transaction
	if err = fn(transaction); err != nil {
		return err
	}

	return transaction.Commit()
}
