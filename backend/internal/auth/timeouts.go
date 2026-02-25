package auth

import "time"

const (
	REGENERATE_TOKEN_TIMEOUT = 10 * time.Second

	SESSION_EXTENSION_DURATION  = 7 * 24 * time.Hour // Extend the session by this time
	SESSION_EXTENSION_CONDITION = 7 * 24 * time.Hour // Extend the session if it expires within this time
)
