package rtoken

import (
	"context"
	"database/sql"
	"time"
)

// TODO: are evita magic link per e 2fa (anche perché 2fa dopo deve generi token di accesso, quindi dev'essere sul dispositivo su cui si vuole accedere)

type IRefreshTokenRepository interface {
	// Create a new token (or one that replaces an old one) and update the struct "refreshToken" with the ExpiresAt and CreatedAt
	CreateToken(ctx context.Context, transaction *sql.Tx, refreshToken *RefreshToken, tokenExp time.Duration) error
	GetToken(ctx context.Context, transaction *sql.Tx, hashedToken []byte) (*RefreshToken, error)
	RevokeTokenById(ctx context.Context, transaction *sql.Tx, tokenId int64, revokedAt time.Time) error
	// DeleteSessionById(ctx context.Context, userId int64, sessionId uuid.UUID) error
}

// DeleteExpired[User]Sessions(ctx context.Context, ...) error

/*
	New session:

	create new token with
	- user_id
	- session_id
	- token_hash
	- expires_at


	Token Rotation:

	(get expired token) -> id, user_id (?), session_id, expires_at, created_at

	transaction {
		create new token with
		- user_id (check with old token?)
		- session_id (from old token)
		- token_hash
		- expires_at (from old token (same as old or with added time))
		- replaces (id of old token)

		update old token with
		- revoked_at (created_at from new token)
	}
*/
