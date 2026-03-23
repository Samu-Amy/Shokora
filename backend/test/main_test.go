package main

import (
	"os"
	"testing"

	"github.com/Samu-Amy/Shokora/internal/api"
	"github.com/Samu-Amy/Shokora/internal/appconfig"
	"github.com/Samu-Amy/Shokora/internal/database"
	"github.com/Samu-Amy/Shokora/internal/service"
	"github.com/Samu-Amy/Shokora/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

var testService *service.Service

var testRouter *chi.Mux

func TestMain(m *testing.M) {

	// - App and DB Config -
	configs := appconfig.NewTestConfig()

	// - Logger -
	logger := zap.NewNop().Sugar()

	// - Mailer -
	mailer := appconfig.GetMailerFromConfig(configs)

	// - Authenticators -
	jwtAuthenticator := appconfig.GetJWTAuthenticatorFromConfig(configs)

	tokenAuthenricator := appconfig.GetTokenAuthenticatorFromConfig(configs)

	// - DB Connection -
	db, err := appconfig.GetDbFromConfig(configs)
	if err != nil {
		panic(err)
	}

	logger.Info("DB Connected")

	// - Transaction Manager -
	txManager := database.NewSQLTransactionManager(db)

	// - Store -
	store := store.NewPostgresStorage(db)

	// - Service -
	authServiceConfig := appconfig.GetAuthServiceConfig(configs)
	testService = service.NewService(txManager, store, mailer, logger, jwtAuthenticator, tokenAuthenricator, authServiceConfig)

	// - Rate Limiter -
	rateLimiter := appconfig.GetFixedWindowLimiterFromConfig(configs)

	// - App -
	testApp := api.NewApp(
		configs,
		testService,
		logger,
		rateLimiter,
	)

	// - Router -
	testRouter = testApp.InitRouter() // Useful for http tests

	code := m.Run()

	// os.Exit skip defer, so we must clean up here
	db.Close()
	logger.Sync()

	os.Exit(code)
}
