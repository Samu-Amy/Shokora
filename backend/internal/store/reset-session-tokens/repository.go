package rstoken

import (
	"context"
)

type IResetSessionTokenRepository interface {
	Create(ctx context.Context, token *RSToken) error

	Get(ctx context.Context, hashedToken []byte) (*RSToken, error)

	Delete(ctx context.Context, hashedToken []byte) error
}
