package auth

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	UserId      int64
	SessionId   uuid.UUID
	HashedToken []byte
	Exp         time.Duration
	Replaces    *int64
}

type RefreshTokenPayload struct {
	PlainToken string
	ExpiresAt  time.Time
}

func GenerateSessionId() (uuid.UUID, error) {
	id, err := uuid.NewRandom()
	return id, err
}
