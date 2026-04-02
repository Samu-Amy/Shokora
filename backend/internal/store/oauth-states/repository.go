package oauthstate

import (
	"context"
	"database/sql"
)

type IOAuthStateRepository interface {
	Create(ctx context.Context, state string) error
	// Get(ctx context.Context, oAuthState *OAuthState) error // Get the state and createdAt
	Delete(ctx context.Context, transaction *sql.Tx, state string) error
}
