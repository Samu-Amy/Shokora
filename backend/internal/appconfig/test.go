package appconfig

import (
	"fmt"
	"time"

	"github.com/Samu-Amy/Shokora/internal/api/ratelimiter"
	"github.com/Samu-Amy/Shokora/internal/config"
	"github.com/Samu-Amy/Shokora/internal/env"
)

func NewTestConfig() config.Config {
	environment := env.GetString("ENV", "test")

	if environment != "prod" {
		env.LoadEnv() // Dev/test Only
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
			Addr:         fmt.Sprintf("host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable", env.GetString("POSTGRES_PORT", "5432"), env.GetString("POSTGRES_TEST_USER", ""), env.GetString("POSTGRES_TEST_PASSWORD", ""), env.GetString("POSTGRES_TEST_DB", "")),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30), // TODO: usare questi valori o lasciare quelli di base?
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		Mail: config.MailConfig{
			Resend: config.ResendConfig{
				ApiKey: env.GetString("RESEND_API_KEY", ""),
			},
			FromEmail:    env.GetString("FROM_EMAIL", ""),
			IsSandboxEnv: env.GetBool("SANDBOX_TEST", true),
		},
		Auth: config.AuthConfig{
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
				Length:      6, // Suggested: between 4 and 10
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
