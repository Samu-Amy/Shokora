package appconfig

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Samu-Amy/Shokora/internal/api/ratelimiter"
	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/config"
	"github.com/Samu-Amy/Shokora/internal/database"
	"github.com/Samu-Amy/Shokora/internal/env"
	"github.com/Samu-Amy/Shokora/internal/mailer"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	OtpLength uint8 = 6
)

func NewDefaultConfig() config.Config {
	environment := env.GetString("ENV", "dev")

	if environment != "prod" {
		env.LoadDevEnv() // Dev Only
	}

	return config.Config{
		Addr:        env.GetString("SERVER_PORT", ":8080"),
		Env:         environment,
		FrontEndURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		AllowedOriginsURLs: env.LoadCORSOrigins([]string{
			"http://localhost:5173",
			"http://localhost:3000",
		}),
		Db: config.DbConfig{
			// Addr: fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=%s", env.GetString("POSTGRES_USER", "user"), env.GetString("POSTGRES_PASSWORD", "password"), env.GetString("POSTGRES_DB", "db"), env.GetString("POSTGRES_PORT", "5432"), env.GetString("POSTGRES_SSL_MODE", "disable")),
			// TODO: attivare modalità ssl (?)
			Addr:         fmt.Sprintf("host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable", env.GetString("POSTGRES_PORT", "5432"), env.GetString("POSTGRES_USER", ""), env.GetString("POSTGRES_PASSWORD", ""), env.GetString("POSTGRES_DB", "")),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30), // TODO: usare questi valori o lasciare quelli di base?
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		Mail: config.MailConfig{
			Resend: config.ResendConfig{
				ApiKey: env.GetString("RESEND_API_KEY", ""),
			},
			FromEmail:    env.GetString("FROM_EMAIL", ""),
			IsSandboxEnv: env.GetBool("SANDBOX", false),
		},
		Auth: config.AuthConfig{
			GoogleOAuthConfig: &oauth2.Config{
				ClientID:     env.GetString("GOOGLE_CLIENT_ID", ""),
				ClientSecret: env.GetString("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  env.GetString("GOOGLE_REDIRECT_URL", ""),
				Scopes: []string{
					"openid",
					"email",
					"profile",
				},
				Endpoint: google.Endpoint,
			},
			PasswordHashingCost: 12, // bcrypt.DefaultCost = 10
			Token: config.TokensConfig{
				Secret:                    env.GetString("AUTH_TOKEN_SECRET", "4a4c345b5064c9a85fff749313ff25310085a606a47232c94e9d898470c6e854"), // TODO: cambia quello di default (oppure dai errore se non lo trova dall'env) e usa "openssl rand -hex 32" per generare token
				Audience:                  "shokora",
				Issuer:                    "shokora",
				ResetSessionTokenByteSize: 32,
				ResetSessionTokenExp:      10 * time.Minute,
				AccessTokenExp:            15 * time.Minute, // 15 min (suggested: 15-60 min) //TODO: alza a 30 (?)
				RefreshTokenByteSize:      32,
				RefreshTokenExp:           30 * 24 * time.Hour, // 30 days (suggested: 7-30 days)
				SessionExp:                90 * 24 * time.Hour, // 90 days (suggested: max 90 days)
			},
			MagicLink: config.MagicLinkConfig{
				ByteSize: 32,
				Exp:      30 * time.Minute, // 30 min
			},
			OTP: config.OTPConfig{
				Length:      OtpLength, // Suggested: between 4 and 10
				MaxAttempts: 5,
				LongExp:     10 * time.Minute, // 10 min (email verification, password reset)
				BaseExp:     5 * time.Minute,  // 5 min (2FA)
			},
			VerficationTokensSecret: env.GetString("VERIFICATION_TOKENS_SECRET", "076477e061001e898408230972c4ec67b806b38449860c8304e04e0ef33b60be"),
		},
		RateLimiter: ratelimiter.RateLimiterConfig{
			RequestsPerTimeFrame: env.GetInt("RATE_LIMITER_REQUESTS_COUNT", 20), // TODO: fix (cambia)
			TimeFrame:            5 * time.Second,
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", true), // TODO: sistema rate limiter e riattivalo
		},
	}
}

func GetMailerFromConfig(configs config.Config) *mailer.ResendMailer {
	return mailer.NewResendMailer(configs.Mail.Resend.ApiKey, configs.Mail.FromEmail)
}

func GetJWTAuthenticatorFromConfig(configs config.Config) *auth.JWTAuthenticator {
	return auth.NewJWTAuthenticator(
		configs.Auth.Token.Secret,
		configs.Auth.Token.Audience,
		configs.Auth.Token.Issuer,
	)
}

func GetTokenAuthenticatorFromConfig(configs config.Config) *auth.TokenAuthenticator {
	return auth.NewTokenAuthenticator(
		configs.Auth.MagicLink,
		configs.Auth.OTP,
		configs.Auth.VerficationTokensMaxRetries,
		configs.Auth.VerficationTokensSecret,
	)
}

func GetDbFromConfig(configs config.Config) (*sql.DB, error) {
	return database.New(
		configs.Db.Addr,
		configs.Db.MaxOpenConns,
		configs.Db.MaxIdleConns,
		configs.Db.MaxIdleTime,
		true,
	)
}

func GetAuthServiceConfig(configs config.Config) config.AuthServiceConfig {
	return config.AuthServiceConfig{
		PasswordHashingCost: configs.Auth.PasswordHashingCost,
		Token:               configs.Auth.Token,
		Mail: config.MailerConfig{
			FrontEndURL:  configs.FrontEndURL,
			FromEmail:    configs.Mail.FromEmail,
			IsSandboxEnv: configs.Mail.IsSandboxEnv,
		},
		Auth: configs.Auth,
	}
}

func GetFixedWindowLimiterFromConfig(configs config.Config) *ratelimiter.FixedWindowLimiter {
	return ratelimiter.NewFixedWindowLimiter(
		configs.RateLimiter.RequestsPerTimeFrame,
		configs.RateLimiter.TimeFrame,
	)
}
