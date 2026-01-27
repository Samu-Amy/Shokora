package store

import (
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
)

type VerificationTokens struct {
	UserId             int64                 `json:"user_id"`
	VerificationType   auth.VerificationType `json:"verification_type"`
	MagicLinkTokenHash []byte                `json:"-"` // TODO: va bene "-"?
	MagicLinkTokenExp  time.Time             `json:"magic_link_token_exp"`
	OTPHash            []byte                `json:"-"`
	OTPExp             time.Time             `json:"otp_exp"`
	OTPAttempts        uint8                 `json:"otp_attempts"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at"`
}
