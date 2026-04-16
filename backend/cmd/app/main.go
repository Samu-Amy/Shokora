package main

import (
	"expvar"
	"runtime"

	"github.com/Samu-Amy/Shokora/internal/api"
	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/appconfig"
	"github.com/Samu-Amy/Shokora/internal/database"
	"github.com/Samu-Amy/Shokora/internal/service"
	"github.com/Samu-Amy/Shokora/internal/store"
	"go.uber.org/zap"
)

// TODO: JWT in HTTP only cookies (no in local storage per evitare XSS) -> attenzione a CSRF (cross origin requests)

// TODO: fai test con/senza redis (sia con dati in cache che non in cache) calcolando il tempo impiegato (?)

// DB Connection string

func main() {

	// - App and DB Config -
	config := appconfig.NewDefaultConfig()

	// - Logger -
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// - Mailer -
	mailer := appconfig.GetMailerFromConfig(config)

	// - Authenticators -
	jwtAuthenticator := appconfig.GetJWTAuthenticatorFromConfig(config)

	tokenAuthenricator := appconfig.GetTokenAuthenticatorFromConfig(config)

	// - DB Connection -
	db, err := appconfig.GetDbFromConfig(config)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("DB Connected")

	// - Transaction Manager -
	txManager := database.NewSQLTransactionManager(db)

	// - Store -
	store := store.NewPostgresStorage(db)

	// - Validator -
	dataValidator := payloads.NewValidator()

	// - Service -
	authServiceConfig := appconfig.GetAuthServiceConfig(config)
	service := service.NewService(txManager, store, mailer, logger, jwtAuthenticator, tokenAuthenricator, authServiceConfig)

	// - Rate Limiter -
	rateLimiter := appconfig.GetFixedWindowLimiterFromConfig(config)

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
		config,
		dataValidator,
		service,
		logger,
		rateLimiter,
	)

	err = app.Run()

	if err != nil {
		logger.Error(err)
	}
}
