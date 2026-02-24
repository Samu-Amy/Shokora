package auth

import (
	"time"
)

type CreateRefreshTokenPayload struct {
	PlainToken string
	ExpiresAt  time.Time
}
