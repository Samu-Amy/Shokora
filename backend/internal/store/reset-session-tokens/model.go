package rstokens

import "time"

// Reset Session Token
type RSToken struct {
	TokenHash []byte
	UserId    int64
	ExpiresAt time.Time
	CreatedAt time.Time
}
