package rtoken

import "time"

type TokenAndSessionData struct {
	Id               int64
	SessionId        int64
	ExpiresAt        time.Time
	RevokedAt        *time.Time
	UserId           int64
	SessionExpiresAt time.Time
}
