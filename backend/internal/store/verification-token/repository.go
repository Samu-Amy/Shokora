package vtoken

import (
	"context"
	"database/sql"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/google/uuid"
)

type IVTokenRepository interface {
	// Create (or update, if already exist) magic link token and otp for email verification | password reset | 2FA
	Create(ctx context.Context, userId int64, verificationTokens *auth.VerificationTokens) error

	UpdateMagicLinkTokenFromId(ctx context.Context, verificationId uuid.UUID, magicLinkTokenHash []byte, magicLinkTokenExp time.Duration) error
	UpdateOTPFromId(ctx context.Context, verificationId uuid.UUID, otpHash []byte, otpExp time.Duration) error

	IncrementOtpAttempts(ctx context.Context, transaction *sql.Tx, verificationId uuid.UUID, maxOTPAttempts uint8) error

	GetOtpData(ctx context.Context, transaction *sql.Tx, verificationId uuid.UUID, verificationType auth.VerificationType) (*OTPVerificationData, error)

	// Get MagicLinkTokenPayload if magic link token found and is not expired
	GetValidMagicLinkData(ctx context.Context, transaction *sql.Tx, hashedToken []byte, verificationType auth.VerificationType) (*MagicLinkVerificationData, error)

	Delete(ctx context.Context, transaction *sql.Tx, verificationId uuid.UUID) error
}
