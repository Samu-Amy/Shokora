package oauthstate

import (
	"time"
)

type OAuthState struct {
	State     string
	CreatedAt time.Time
}
