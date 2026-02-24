package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// TODO: modifica (ora anche tabella sessione)
type RefreshToken struct {
	Id        int64      `json:"id"` // Generated
	SessionId int64      `json:"session_id"`
	TokenHash []byte     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	Replaces  *int64     `json:"replaces,omitempty"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"` // Default now()
}

// TODO: are evita magic link per e 2fa (anche perché 2fa dopo deve generi token di accesso, quindi dev'essere sul dispositivo su cui si vuole accedere)

// Repository
type RefreshTokenRepositoryI interface {
	// Create a new token (or one that replaces an old one) and update the struct "refreshToken" with the ExpiresAt and CreatedAt
	CreateToken(ctx context.Context, transaction *sql.Tx, refreshToken *RefreshToken, tokenExp time.Duration) error

	GetToken(ctx context.Context, transaction *sql.Tx, hashedToken []byte) (*RefreshToken, error)

	RevokeTokenById(ctx context.Context, transaction *sql.Tx, tokenId int64, revokedAt time.Time) error

	DeleteSessionById(ctx context.Context, userId int64, sessionId uuid.UUID) error

	// TODO: nel login fai anche delete di tutti i refresh token scaduti per quell'utente (o in generale?) - ottenere un l'ultimo token creato per ogni session_id (join con order by) e se è scaduto -> sessione scaduta (?)
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
}
