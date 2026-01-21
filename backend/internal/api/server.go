package api

import (
	"net/http"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/mailer"
	"github.com/Samu-Amy/Shokora/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// - Structs -
type App struct {
	config        Config
	router        *chi.Mux
	store         *store.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
}

type Config struct {
	Addr        string
	Env         string // "env" | "prod"
	FrontEndURL string
	Db          DbConfig
	Mail        MailConfig
	Auth        AuthConfig
}

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type MailConfig struct {
	Resend                    ResendConfig
	FromEmail                 string
	EmailVerificationTokenExp time.Duration
	PasswordResetTokenExp     time.Duration
}

type ResendConfig struct {
	ApiKey string
}

type AuthConfig struct {
	Token TokenConfig
}

type TokenConfig struct {
	Secret   string
	Audience string
	Issuer   string
	Exp      time.Duration
}

// - Functions/Methods -
func NewApp(
	config Config,
	store *store.Storage,
	logger *zap.SugaredLogger,
	mailer mailer.Client,
	authenticator auth.Authenticator,
) *App {
	app := &App{
		config:        config,
		store:         store,
		logger:        logger,
		mailer:        mailer,
		authenticator: authenticator,
	}
	app.router = app.initRouter()

	return app
}

func (app *App) Run() error {
	server := &http.Server{
		Addr:         app.config.Addr,
		Handler:      app.router,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("Server started", "addr", app.config.Addr)
	return server.ListenAndServe()
}

// - Tests -

func NewMockApp(store *store.Storage, logger *zap.SugaredLogger, authenticator auth.Authenticator) *App {
	app := &App{
		store:         store,
		logger:        logger,
		authenticator: authenticator,
	}
	app.router = app.initRouter()

	return app
}

func (app *App) GetRouter() *chi.Mux {
	return app.router
}

func (app *App) GetAuthenticator() auth.Authenticator {
	return app.authenticator
}
