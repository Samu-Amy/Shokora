package store

import (
	"context"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/google/uuid"
)

type RefreshTokens struct {
	Id        int64      `json:"id"` // Generated
	UserId    int64      `json:"user_id"`
	SessionId uuid.UUID  `json:"session_id"`
	TokenHash []byte     `json:"-"`
	Exp       time.Time  `json:"expires_at"`
	Replaces  *int64     `json:"replaces,omitempty"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"` // Default now()
}

// TODO: evita magic link per e 2fa (anche perché 2fa dopo deve generare i token di accesso, quindi dev'essere sul dispositivo su cui si vuole accedere)

// Repository
type RefreshTokensRepositoryI interface {
	CreateToken(ctx context.Context, refreshToken auth.RefreshToken) error

	// GetToken(ctx context.Context, hashedToken []byte) (..., error) // TODO: ritorna dati che servono (es. session_id, created_at (per revoked)) - usa in service (fai tutto in transaction)

	RevokeTokenById(ctx context.Context, tokenId int64, revokedAt time.Time) error

	// TODO: nel login fai anche delete di tutti i refresh token scaduti per quell'utente (?)

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
