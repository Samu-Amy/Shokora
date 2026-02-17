package auth

import (
	"time"

	"github.com/google/uuid"
)

// Used to create a new token or a one that replaces another (in -> data into db, out -> data from the db)
type RefreshToken struct { // TODO: sostituire con RefreshTokens (store)
	UserId      int64
	SessionId   uuid.UUID
	HashedToken []byte
	Exp         time.Duration
	Replaces    *int64     // in
	ExpiresAt   *time.Time // out
	CreatedAt   *time.Time // out
}

type CreateRefreshTokenPayload struct {
	PlainToken string
	ExpiresAt  time.Time
}

func GenerateSessionId() (uuid.UUID, error) {
	id, err := uuid.NewRandom()
	return id, err
}
