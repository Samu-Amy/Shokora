package rtoken

import (
	"time"
)

type CreateRefreshTokenDto struct {
	PlainToken string
	ExpiresAt  time.Time
}
