package vtoken

import (
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
)

// Verification Tokens (Magic Link and OTP)
type VToken struct {
	Id                      int64
	UserId                  int64
	VerificationType        auth.VerificationType
	MagicLinkTokenHash      []byte
	MagicLinkTokenExpiresAt time.Time
	OTPHash                 []byte
	OTPExpiresAt            time.Time
	OTPAttempts             uint8
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

// - Payloads -

type MagicLinkTokenPayload struct {
	VerificationId int64
	UserId         int64
}

type OTPPayload struct {
	UserId    int64
	HashedOtp []byte
	Attempts  uint8
	ExpiresAt time.Time
}
