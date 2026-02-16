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

func GenerateSessionId() (*uuid.UUID, error) {
	return nil, nil
}
