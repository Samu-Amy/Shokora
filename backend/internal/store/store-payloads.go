package store

import (
	"time"
)

// - Payloads -

type MagicLinkTokenPayload struct {
	VerificationId int64
	UserId         int64
}

type OTPPayload struct {
	UserId    int64
	HashedOtp []byte
	Attempts  uint8
	Exp       time.Time
}
