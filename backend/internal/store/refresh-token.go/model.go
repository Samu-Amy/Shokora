package rtoken

import (
	"time"
)

type RefreshToken struct {
	Id        int64
	SessionId int64
	TokenHash []byte
	ExpiresAt time.Time
	Replaces  *int64
	RevokedAt *time.Time
	CreatedAt time.Time
}
