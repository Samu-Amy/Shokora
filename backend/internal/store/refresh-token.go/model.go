package rtoken

import (
	"time"
)

// TODO: modifica (ora anche tabella sessione)
type RefreshToken struct {
	Id        int64
	SessionId int64
	TokenHash []byte
	ExpiresAt time.Time
	Replaces  *int64
	RevokedAt *time.Time
	CreatedAt time.Time
}
