package config

import (
	"time"

	"github.com/Samu-Amy/Shokora/internal/api/ratelimiter"
	"golang.org/x/oauth2"
)

// ----- APP -----

type Config struct {
	Addr               string
	Env                string // "env" | "prod" | "test"
	FrontEndURL        string
	AllowedOriginsURLs []string
	Db                 DbConfig
	Mail               MailConfig
	Auth               AuthConfig
	RateLimiter        ratelimiter.RateLimiterConfig
}

// - Services -

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type MailConfig struct {
	Resend       ResendConfig
	FromEmail    string
	IsSandboxEnv bool
}

type ResendConfig struct {
	ApiKey string
}

// ----- AUTH -----

type AuthConfig struct {
	GoogleOAuthConfig           *oauth2.Config
	PasswordHashingCost         int
	Token                       TokensConfig
	MagicLink                   MagicLinkConfig
	OTP                         OTPConfig
	VerficationTokensMaxRetries uint8 // Counting the first attempt
	VerficationTokensSecret     string
}

type TokensConfig struct {
	Secret                    string
	Audience                  string
	Issuer                    string
	ResetSessionTokenByteSize int
	ResetSessionTokenExp      time.Duration
	AccessTokenExp            time.Duration
	RefreshTokenByteSize      int
	RefreshTokenExp           time.Duration
	SessionExp                time.Duration // How long the refresh tokens expiration can be extended for
}

// - Verification -

type MagicLinkConfig struct {
	ByteSize int
	Exp      time.Duration
}

type OTPConfig struct {
	Length      uint8
	MaxAttempts uint8
	LongExp     time.Duration // For Email Verification
	BaseExp     time.Duration // For Password Reset and 2FA
	// CriticalExp time.Duration // For critical operations (es. 30s)
}

// ----- SERVICE LAYER ----

// - Auth -

type AuthServiceConfig struct {
	PasswordHashingCost int
	Token               TokensConfig
	Mail                MailerConfig
	Auth                AuthConfig
}

type MailerConfig struct {
	FrontEndURL  string
	FromEmail    string
	IsSandboxEnv bool
}
