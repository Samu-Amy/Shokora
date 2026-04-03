package rstoken

import (
	"context"
	"database/sql"
)

type IResetSessionTokenRepository interface {
	Create(ctx context.Context, transaction *sql.Tx, rsToken *RSToken) error

	Get(ctx context.Context, hashedToken []byte) (*RSToken, error)

	Delete(ctx context.Context, transaction *sql.Tx, hashedToken []byte) error
}
