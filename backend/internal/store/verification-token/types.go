package vtoken

import "time"

type MagicLinkVerificationData struct {
	VerificationId int64
	UserId         int64
}

type OTPVerificationData struct {
	UserId    int64
	HashedOtp []byte
	Attempts  uint8
	ExpiresAt time.Time
}
