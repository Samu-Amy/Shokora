package main

import (
	"fmt"
	"time"

	"github.com/Samu-Amy/Shokora/internal/api"
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
		FrontEndURL: env.GetString("FRONTEND_URL", "http://localhost:3000"),
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

	// - App -
	app := api.NewApp(config, &store, logger, mailer)

	err = app.Run()

	if err != nil {
		logger.Error(err) // TODO: sistema
	}
}
