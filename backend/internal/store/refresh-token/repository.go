package rtoken

import (
	"context"
	"database/sql"
	"time"

	session "github.com/Samu-Amy/Shokora/internal/store/user-session"
)

type IRefreshTokenRepository interface {
	// Create a new refresh token (or one that replaces an old one) and update the struct "refreshToken" with the ExpiresAt and CreatedAt
	Create(ctx context.Context, transaction *sql.Tx, refreshToken *RefreshToken) error

	// Get the refresh token data, userId and session ExpiresAt
	GetByTokenForUpdate(ctx context.Context, transaction *sql.Tx, hashedToken []byte) (*TokenAndSessionData, error)
	GetSessionDataByToken(ctx context.Context, hashedToken []byte) (*session.SessionData, error)

	RevokeById(ctx context.Context, transaction *sql.Tx, tokenId int64, revokedAt time.Time) error
}
