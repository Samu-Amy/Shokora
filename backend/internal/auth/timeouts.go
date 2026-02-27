package auth

import "time"

const (
	RegenerateTokenTimeout = 10 * time.Second

	SessionExtensionDuration  = 7 * 24 * time.Hour // Extend the session by this time
	SessionExtensionCondition = 7 * 24 * time.Hour // Extend the session if it expires within this time
)
