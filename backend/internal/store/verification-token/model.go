package v_token

import (
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
)

// Verification Tokens (Magic Link and OTP)
type VToken struct {
	Id                 int64                 `json:"id"` // Generated
	UserId             int64                 `json:"user_id"`
	VerificationType   auth.VerificationType `json:"verification_type"`
	MagicLinkTokenHash []byte                `json:"-"`
	MagicLinkTokenExp  time.Time             `json:"magic_link_token_exp"`
	OTPHash            []byte                `json:"-"`
	OTPExp             time.Time             `json:"otp_exp"`
	OTPAttempts        uint8                 `json:"otp_attempts"` // Default 0 (otp attempts for (user_id, verificationType))
	CreatedAt          time.Time             `json:"created_at"`   // Default now()
	UpdatedAt          time.Time             `json:"updated_at"`   // Default now()
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
	Exp       time.Time
}
