package store

import (
	"context"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
)

// Verification Tokens (Magic Link and OTP)
type VTokens struct {
	UserId             int64                 `json:"user_id"`
	VerificationType   auth.VerificationType `json:"verification_type"`
	MagicLinkTokenHash []byte                `json:"-"` // TODO: va bene "-"?
	MagicLinkTokenExp  time.Time             `json:"magic_link_token_exp"`
	OTPHash            []byte                `json:"-"`
	OTPExp             time.Time             `json:"otp_exp"`
	OTPAttempts        uint8                 `json:"otp_attempts"` // Default 0
	CreatedAt          time.Time             `json:"created_at"`   // Default now()
	UpdatedAt          time.Time             `json:"updated_at"`   // Default now()
}

// Repository
type VTokensRepositoryI interface {
	// Create tokens (for email verification | password reset | 2FA)
	CreateTokens(ctx context.Context, userId int64, verificationTokens *auth.VerificationTokens) error

	UpdateMagicLinkToken(ctx context.Context, userId int64, magicLinkTokenHash []byte, magicLinkTokenExp time.Duration) error
	UpdateOTP(ctx context.Context, userId int64, OTPHash []byte, OTPExp time.Duration) error
}
