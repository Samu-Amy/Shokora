package store

import (
	"context"
	"errors"
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

// Errors
var (
	ErrDuplicateMagicLinkToken = errors.New("duplicate_token") // TODO: implementa errori
	ErrDuplicateOTP            = errors.New("duplicate_otp")
)

// Repository
type VTokensRepositoryI interface {
	// Create tokens (for email verification | password reset | 2FA)
	CreateTokens(ctx context.Context, verificationTokens *auth.VerificationTokens, userId int64) error
}
