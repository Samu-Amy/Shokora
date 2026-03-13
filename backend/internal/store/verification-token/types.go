package vtoken

import (
	"time"

	"github.com/google/uuid"
)

type MagicLinkVerificationData struct {
	VerificationId uuid.UUID
	UserId         int64
}

type OTPVerificationData struct {
	UserId    int64
	HashedOtp []byte
	Attempts  uint8
	ExpiresAt time.Time
}
