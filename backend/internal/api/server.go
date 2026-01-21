package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Samu-Amy/Shokora/internal/api/ratelimiter"
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
	rateLimiter   ratelimiter.RateLimiter
}

type Config struct {
	Addr        string
	Env         string // "env" | "prod"
	FrontEndURL string
	Db          DbConfig
	Mail        MailConfig
	Auth        AuthConfig
	RateLimiter ratelimiter.RateLimiterConfig
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
	rateLimiter ratelimiter.RateLimiter,
) *App {
	app := &App{
		config:        config,
		store:         store,
		logger:        logger,
		mailer:        mailer,
		authenticator: authenticator,
		rateLimiter:   rateLimiter,
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

	// Graceful shutdown
	shutdownErr := make(chan error, 1)

	// TODO: sistemare (https://www.youtube.com/watch?v=UPVSeZXBTxI) (?)
	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		// defer signal.Stop(quit)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdownErr <- server.Shutdown(ctx)
	}()

	app.logger.Infow("Server started", "addr", app.config.Addr)

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErr
	if err != nil {
		return err
	}

	app.logger.Warnf("server has stopped", "addr", app.config.Addr, "env", app.config.Env)

	return nil
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
