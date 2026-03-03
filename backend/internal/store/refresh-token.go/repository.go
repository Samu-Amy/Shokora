package rtoken

import (
	"context"
	"database/sql"
	"time"
)

type IRefreshTokenRepository interface {
	// Create a new refresh token (or one that replaces an old one) and update the struct "refreshToken" with the ExpiresAt and CreatedAt
	Create(ctx context.Context, transaction *sql.Tx, refreshToken *RefreshToken) error

	// Get the refresh token data, userId and session ExpiresAt
	GetByToken(ctx context.Context, transaction *sql.Tx, hashedToken []byte) (*TokenAndSessionData, error)
	RevokeById(ctx context.Context, transaction *sql.Tx, tokenId int64, revokedAt time.Time) error
}
