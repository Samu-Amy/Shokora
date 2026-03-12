package vtoken

import (
	"context"
	"database/sql"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
)

// TODO: evita magic link per e 2fa (anche perché 2fa dopo deve generare i token di accesso, quindi dev'essere sul dispositivo su cui si vuole accedere)

// TODO: usare VToken invece di auth.VerificationTokens con le exp (Duration) separate e ritornare i dati in VToken settando i campi non inizializzati (come faccio in User e RefreshToken)

type IVTokenRepository interface {
	// Create (or update, if already exist) magic link token and otp for email verification | password reset | 2FA
	Create(ctx context.Context, userId int64, verificationTokens *auth.VerificationTokens) error

	UpdateMagicLinkTokenFromId(ctx context.Context, verificationId int64, magicLinkTokenHash []byte, magicLinkTokenExp time.Duration) error
	UpdateOTPFromId(ctx context.Context, verificationId int64, otpHash []byte, otpExp time.Duration) error

	UpdateOtpAttempts(ctx context.Context, transaction *sql.Tx, verificationId int64, maxOTPAttempts uint8) error

	GetOtpData(ctx context.Context, transaction *sql.Tx, verificationId int64, verificationType auth.VerificationType) (*OTPVerificationData, error)

	// Get MagicLinkTokenPayload if magic link token found and is not expired
	GetValidMagicLinkData(ctx context.Context, hashedToken []byte, verificationType auth.VerificationType) (*MagicLinkVerificationData, error)

	Delete(ctx context.Context, verificationId int64) error
}
