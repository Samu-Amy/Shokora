package store

import "time"

type UserSession struct {
	Id        int64     `json:"id"` // Generated
	UserId    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"` // Default now()
}
