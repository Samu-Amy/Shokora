package database

import (
	"context"
	"database/sql"
)

type ITransactionManager interface {
	WithTx(ctx context.Context, fn func(*sql.Tx) error) error
}
