package oauthstate

import (
	"context"
)

type IOAuthStateRepository interface {
	Create(ctx context.Context, state string) error
	Get(ctx context.Context, oAuthState *OAuthState) error // Get the state and createdAt
	Delete(ctx context.Context, state string) error
}
