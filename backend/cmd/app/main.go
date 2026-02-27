package main

import (
	"expvar"
	"fmt"
	"runtime"
	"time"

	"github.com/Samu-Amy/Shokora/internal/api"
	"github.com/Samu-Amy/Shokora/internal/api/ratelimiter"
	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/config"
	"github.com/Samu-Amy/Shokora/internal/database"
	"github.com/Samu-Amy/Shokora/internal/env"
	"github.com/Samu-Amy/Shokora/internal/mailer"
	"github.com/Samu-Amy/Shokora/internal/service"
	"github.com/Samu-Amy/Shokora/internal/store"
	"go.uber.org/zap"
)

// TODO: JWT in HTTP only cookies (no in local storage per evitare XSS) -> attenzione a CSRF (cross origin requests)

// TODO: fai test con/senza redis (sia con dati in cache che non in cache) calcolando il tempo impiegato (?)

// DB Connection string

func main() {
	environment := env.GetString("ENV", "dev")

	if environment != "prod" {
		env.LoadEnv() //! - Dev Only - !
	}

	// - App and DB Config -
	configs := config.Config{
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
			PasswordHashingCost: 12, // bcrypt.DefaultCost = 10
			Token: config.TokenConfig{
				Secret:               env.GetString("AUTH_TOKEN_SECRET", "4a4c345b5064c9a85fff749313ff25310085a606a47232c94e9d898470c6e854"), // TODO: cambia quello di default (oppure dai errore se non lo trova dall'env) e usa "openssl rand -hex 32" per generare token
				Audience:             "shokora",
				Issuer:               "shokora",
				AccessTokenExp:       15 * time.Minute, // 15 min (suggested: 15-60 min) //TODO: alza a 30 (?)
				RefreshTokenByteSize: 32,
				RefreshTokenExp:      30 * 24 * time.Hour, // 30 days (suggested: 7-30 days)
				SessionMaxExp:        90 * 24 * time.Hour, // 90 days (suggested: max 90 days)
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

	// - Logger -
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// - Mailer -
	mailer := mailer.NewResendMailer(configs.Mail.Resend.ApiKey, configs.Mail.FromEmail)

	// - Authenticators -
	jwtAuthenticator := auth.NewJWTAuthenticator(
		configs.Auth.Token.Secret,
		configs.Auth.Token.Audience,
		configs.Auth.Token.Issuer,
	)

	tokenAuthenricator := auth.NewTokenAuthenticator(
		configs.Auth.MagicLink,
		configs.Auth.OTP,
		configs.Auth.VerficationTokensMaxRetries,
		configs.Auth.VerficationTokensSecret,
	)

	// - DB Connection -
	db, err := database.New(
		configs.Db.Addr,
		configs.Db.MaxOpenConns,
		configs.Db.MaxIdleConns,
		configs.Db.MaxIdleTime,
		true,
	)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("DB Connected")

	// - Transaction Manager -
	txManager := database.NewSQLTransactionManager(db)

	// - Store -
	store := store.NewPostgresStorage(db)

	// - Service -
	authServiceConfig := config.AuthServiceConfig{
		PasswordHashingCost: configs.Auth.PasswordHashingCost,
		Token:               configs.Auth.Token,
		Mail: config.MailerConfig{
			FrontEndURL:  configs.FrontEndURL,
			FromEmail:    configs.Mail.FromEmail,
			IsSandboxEnv: configs.Mail.IsSandboxEnv,
		},
	}
	service := service.NewService(txManager, store, mailer, logger, jwtAuthenticator, tokenAuthenricator, authServiceConfig)

	// - Rate Limiter -
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		configs.RateLimiter.RequestsPerTimeFrame,
		configs.RateLimiter.TimeFrame,
	)

	// - Metrics -
	// Version
	// expvar.NewString("backend_version").Set("1.0")

	// DB Stats
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))

	// Goroutines
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	// - App -
	app := api.NewApp(
		configs,
		service,
		logger,
		rateLimiter,
	)

	err = app.Run()

	if err != nil {
		logger.Error(err)
	}
}
