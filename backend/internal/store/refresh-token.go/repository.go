package rtoken

import (
	"context"
	"database/sql"
	"time"
)

type IRefreshTokenRepository interface {
	// Create a new token (or one that replaces an old one) and update the struct "refreshToken" with the ExpiresAt and CreatedAt
	CreateToken(ctx context.Context, transaction *sql.Tx, refreshToken *RefreshToken, tokenExp time.Duration) error
	GetToken(ctx context.Context, transaction *sql.Tx, hashedToken []byte) (*RefreshToken, error)
	RevokeTokenById(ctx context.Context, transaction *sql.Tx, tokenId int64, revokedAt time.Time) error
	// DeleteSessionById(ctx context.Context, userId int64, sessionId uuid.UUID) error
}
