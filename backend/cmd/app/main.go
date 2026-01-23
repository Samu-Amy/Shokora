package main

import (
	"expvar"
	"fmt"
	"runtime"
	"time"

	"github.com/Samu-Amy/Shokora/internal/api"
	"github.com/Samu-Amy/Shokora/internal/api/ratelimiter"
	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/db"
	"github.com/Samu-Amy/Shokora/internal/env"
	"github.com/Samu-Amy/Shokora/internal/mailer"
	"github.com/Samu-Amy/Shokora/internal/store"
	"go.uber.org/zap"
)

// TODO: JWT in HTTP only cookies (no in local storage per evitare XSS) -> attenzione a CSRF (cross origin requests)

// TODO: fai test con/senza redis (sia con dati in cache che non in cache) calcolando il tempo impiegato (?)

// DB Connection string

func main() {
	env.LoadEnv() //! - Dev Only (use file .env) - !

	// - App and DB Config -
	config := api.Config{
		Addr:        env.GetString("SERVER_PORT", ":8080"),
		Env:         env.GetString("ENV", "dev"),
		FrontEndURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		AllowedOriginsURLs: env.LoadCORSOrigins([]string{
			"http://localhost:5173",
			"http://localhost:3000",
		}),
		Db: api.DbConfig{
			// Addr: fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=%s", env.GetString("POSTGRES_USER", "user"), env.GetString("POSTGRES_PASSWORD", "password"), env.GetString("POSTGRES_DB", "db"), env.GetString("POSTGRES_PORT", "5432"), env.GetString("POSTGRES_SSL_MODE", "disable")),
			// TODO: attivare modalità ssl (?)
			Addr:         fmt.Sprintf("host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable", env.GetString("POSTGRES_PORT", "5432"), env.GetString("POSTGRES_USER", ""), env.GetString("POSTGRES_PASSWORD", ""), env.GetString("POSTGRES_DB", "")),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30), // TODO: usare questi valori o lasciare quelli di base?
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		Mail: api.MailConfig{
			Resend: api.ResendConfig{
				ApiKey: env.GetString("RESEND_API_KEY", ""),
			},
			FromEmail:                 env.GetString("FROM_EMAIL", ""),
			EmailVerificationTokenExp: 24 * time.Hour,
			PasswordResetTokenExp:     30 * time.Minute,
		},
		Auth: api.AuthConfig{
			Token: api.TokenConfig{
				Secret:          env.GetString("AUTH_TOKEN_SECRET", "basicTokenSecret"),
				Audience:        "shokora",
				Issuer:          "shokora",
				AccessTokenExp:  15 * time.Minute,
				RefreshTokenExp: 14 * 24 * time.Hour, // 14 days
			},
		},
		RateLimiter: ratelimiter.RateLimiterConfig{
			RequestsPerTimeFrame: env.GetInt("RATE_LIMITER_REQUESTS_COUNT", 20), // TODO: fix (cambia)
			TimeFrame:            5 * time.Second,
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", true),
		},
	}

	// - Logger -
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// - DB Connection -
	db, err := db.New(
		config.Db.Addr,
		config.Db.MaxOpenConns,
		config.Db.MaxIdleConns,
		config.Db.MaxIdleTime,
		true,
	)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("DB Connected")

	// - Store -
	store := store.NewPostgresStorage(db)

	// - Mailer -
	mailer := mailer.NewResendMailer(config.Mail.Resend.ApiKey, config.Mail.FromEmail)

	// - Authenticator -
	jwtAuthenticator := auth.NewJWTAuthenticator(
		config.Auth.Token.Secret,
		config.Auth.Token.Issuer,
		config.Auth.Token.Issuer,
	)

	// - Rate Limiter -
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		config.RateLimiter.RequestsPerTimeFrame,
		config.RateLimiter.TimeFrame,
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
	app := api.NewApp(config, &store, logger, mailer, jwtAuthenticator, rateLimiter)

	err = app.Run()

	if err != nil {
		logger.Error(err)
	}
}
