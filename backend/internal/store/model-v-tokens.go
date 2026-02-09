package store

import (
	"context"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
)

// Verification Tokens (Magic Link and OTP)
type VTokens struct {
	Id                 int64                 `json:"id"` // Generated
	UserId             int64                 `json:"user_id"`
	VerificationType   auth.VerificationType `json:"verification_type"`
	MagicLinkTokenHash []byte                `json:"-"`
	MagicLinkTokenExp  time.Time             `json:"magic_link_token_exp"`
	OTPHash            []byte                `json:"-"`
	OTPExp             time.Time             `json:"otp_exp"`
	OTPAttempts        uint8                 `json:"otp_attempts"` // Default 0
	CreatedAt          time.Time             `json:"created_at"`   // Default now()
	UpdatedAt          time.Time             `json:"updated_at"`   // Default now()
}

// TODO: evita magic link per e 2fa (anche perché 2fa dopo deve generare i token di accesso, quindi dev'essere sul dispositivo su cui si vuole accedere)

// Repository
type VTokensRepositoryI interface {
	// Create (or update, if already exist) magic link token and otp for email verification | password reset | 2FA
	CreateTokens(ctx context.Context, userId int64, verificationTokens *auth.VerificationTokens) (*int64, error)

	UpdateMagicLinkTokenFromId(ctx context.Context, verificationId int64, magicLinkTokenHash []byte, magicLinkTokenExp time.Duration) error
	UpdateOTPFromId(ctx context.Context, verificationId int64, otpHash []byte, otpExp time.Duration) error

	// TODO: fare metodi apposta per password e/o 2fa (es. nel caso dovessi fare join con users)
	VerifyMagicLink(ctx context.Context, hashedToken []byte) (*int64, *auth.VerificationType, error)
	VerifyOTP(ctx context.Context, verificationI int64, hashedOTP []byte, maxOTPAttempts uint8) (*int64, *auth.VerificationType, error)
}
