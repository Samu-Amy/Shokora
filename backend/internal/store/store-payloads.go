package store

import (
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
)

// - Payloads -

type MagicLinkTokenPayload struct {
	VerificationId   int64
	UserId           int64
	VerificationType auth.VerificationType
	Exp              time.Time
}

type OTPPayload struct {
	UserId           int64
	VerificationType auth.VerificationType
	Attempts         uint8
	Exp              time.Time
}
