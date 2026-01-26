package store

import "context"

func (store *PostgresUserStore) CreateRefreshToken(ctx context.Context, hashedToken []byte) error {
	return nil
}

func (store *PostgresUserStore) RefreshToken(ctx context.Context, hashedToken []byte) error {
	// TODO: usare UNIQUE(token) quando lo si cerca (?)
	return nil
}
