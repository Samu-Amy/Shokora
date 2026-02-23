package auth

import (
	"time"

	"github.com/google/uuid"
)

// Used to create a new token or a one that replaces another (in -> data into db, out -> data from the db)
type RefreshToken struct {
	Id          *int64
	UserId      int64
	SessionId   uuid.UUID
	HashedToken []byte
	Exp         time.Duration // the expiration duration for the token (from config), used in db to set expires_at
	Replaces    *int64        // in
	RevokedAt   *time.Time    // out
	ExpiresAt   *time.Time    // out (expiration date, correspond to the "expires_at" set in db (using "Exp"))
	CreatedAt   *time.Time    // out
}

type CreateRefreshTokenPayload struct {
	PlainToken string
	ExpiresAt  time.Time
}

func GenerateSessionId() (uuid.UUID, error) {
	id, err := uuid.NewRandom()
	return id, err
}
