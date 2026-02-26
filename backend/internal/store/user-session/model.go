package session

import (
	"time"
)

type UserSession struct {
	Id        int64
	UserId    int64
	ExpiresAt time.Time
	CreatedAt time.Time
}
