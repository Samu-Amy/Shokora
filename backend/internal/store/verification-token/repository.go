package v_token

import (
	"context"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
)

// TODO: evita magic link per e 2fa (anche perché 2fa dopo deve generare i token di accesso, quindi dev'essere sul dispositivo su cui si vuole accedere)

type VTokenRepositoryI interface {
	// Create (or update, if already exist) magic link token and otp for email verification | password reset | 2FA
	CreateTokens(ctx context.Context, userId int64, verificationTokens *auth.VerificationTokens) (*int64, error)

	UpdateMagicLinkTokenFromId(ctx context.Context, verificationId int64, magicLinkTokenHash []byte, magicLinkTokenExp time.Duration) error
	UpdateOTPFromId(ctx context.Context, verificationId int64, otpHash []byte, otpExp time.Duration) error

	UpdateOtpAttempts(ctx context.Context, verificationId int64, maxOTPAttempts uint8) error

	GetOtpData(ctx context.Context, verificationId int64, verificationType auth.VerificationType) (*OTPPayload, error)

	VerifyMagicLink(ctx context.Context, hashedToken []byte, verificationType auth.VerificationType) (*MagicLinkTokenPayload, error)

	Delete(ctx context.Context, verificationId int64) error
}
