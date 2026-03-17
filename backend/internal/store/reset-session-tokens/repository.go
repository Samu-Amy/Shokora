package rstokens

import (
	"context"
	"time"
)

type IRSTokenRepository interface {
	Create(ctx context.Context, token *RSToken, tokenExp time.Duration) error

	Get(ctx context.Context, hashedToken []byte) (*RSToken, error)

	Delete(ctx context.Context, hashedToken []byte) error
}
