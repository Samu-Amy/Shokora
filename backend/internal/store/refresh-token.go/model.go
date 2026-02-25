package rtoken

import (
	"time"
)

// TODO: modifica (ora anche tabella sessione)
type RefreshToken struct {
	Id        int64      `json:"id"` // Generated
	SessionId int64      `json:"session_id"`
	TokenHash []byte     `json:"-"`
	ExpiresAt time.Time  `json:"expires_at"`
	Replaces  *int64     `json:"replaces,omitempty"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"` // Default now()
}
