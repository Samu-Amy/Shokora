package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
)

type PostgresVTokensStore struct {
	db *sql.DB
}

func NewPostgresVTokenStore(db *sql.DB) *PostgresVTokensStore {
	return &PostgresVTokensStore{db: db}
}

func (store *PostgresVTokensStore) CreateTokens(ctx context.Context, verificationTokens *auth.VerificationTokens, userId int64) error {
	query := `
		INSERT INTO verification_tokens (token, user_id, expiry)
		VALUES ($1, $2, $3)
	` // TODO: usa UPSERT

	ctx, cancel := context.WithTimeout(ctx, medium_query_timeout)
	defer cancel()

	_, err := store.db.ExecContext(
		ctx,
		query,
		verificationTokens.HashedMagicLinkToken,
		userId,
		time.Now().Add(verificationTokens.OTPExp),
	)

	if err != nil {
		return err
	}

	return nil
}
