package main

import (
	"context"
	"database/sql"
	"math/rand"
	"os"
	"testing"

	"github.com/Samu-Amy/Shokora/internal/api"
	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/appconfig"
	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/database"
	"github.com/Samu-Amy/Shokora/internal/service"
	"github.com/Samu-Amy/Shokora/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

/*
Per runnare i test:
- "cd .\backend\"
	- test normali (tutti): "go test [-v] .\test\"
	- test normali (singolo test): "go test [-v] .\test\ -run [NomeTest (es. TestRegisterUserRoute)]"
	- solo test fuzz: "go test .\test\ -run=^$ -fuzz=[NomeTest (es. FuzzRegisterUserRoute) -fuzztime=[tempo (es. 20s)]" (-run=^$ dice di runnare i test normali che matchano con la regex (nessuno))
*/

// Constants
const (
	activateLogger = true //* Useful for debugging when test don't pass

	routesTestsNum     = 25
	validationTestsNum = 50

	randSeed int64 = 42 //* cambia il seed per testare diversi casi

	// DB Seeding
	seedUserNum = 15
)

// Services to use in test
var customRand *rand.Rand //* usa questo per generare valori pseudo-casuali (con casualità ma anche riproducibilità)
var dataValidator *validator.Validate
var testStore *store.Storage
var db *sql.DB
var testJwtAuthenticator *auth.JWTAuthenticator
var testService *service.Service
var testRouter *chi.Mux

var authState *AuthState

var err error

// Main
func TestMain(m *testing.M) {

	// Set rand seeed
	customRand = rand.New(rand.NewSource(randSeed))

	// - App and DB Config -
	configs := appconfig.NewTestConfig()

	// - Logger -
	var logger *zap.SugaredLogger
	if activateLogger {
		logger = zap.Must(zap.NewProduction()).Sugar()
	} else {
		logger = zap.NewNop().Sugar()
	}

	// - Mailer -
	mailer := appconfig.GetMailerFromConfig(configs)

	// - Authenticators -
	testJwtAuthenticator = appconfig.GetJWTAuthenticatorFromConfig(configs)

	tokenAuthenricator := appconfig.GetTokenAuthenticatorFromConfig(configs)

	// - DB Connection -
	db, err = appconfig.GetDbFromConfig(configs)
	if err != nil {
		panic(err)
	}

	logger.Info("DB Connected")

	// - Transaction Manager -
	txManager := database.NewSQLTransactionManager(db)

	// - Store -
	testStore = store.NewPostgresStorage(db)

	// - Validator -
	dataValidator = payloads.NewValidator()

	// - Service -
	authServiceConfig := appconfig.GetAuthServiceConfig(configs)
	testService = service.NewService(txManager, testStore, mailer, logger, testJwtAuthenticator, tokenAuthenricator, authServiceConfig)

	// - Rate Limiter -
	rateLimiter := appconfig.GetFixedWindowLimiterFromConfig(configs)

	// - App -
	testApp := api.NewApp(
		configs,
		dataValidator,
		testService,
		logger,
		rateLimiter,
	)

	// - Router -
	testRouter = testApp.InitRouter() // Useful for http tests

	// - DB Actions -
	clearTestDB(db)
	authState, err = seedAuthState(context.Background(), db)
	if err != nil {
		logger.Warnf("Error seeding db: %v", err)
		panic(err)
	}

	code := m.Run()

	// os.Exit skip defer, so we must clean up here
	db.Close()
	logger.Sync()

	os.Exit(code)
}
