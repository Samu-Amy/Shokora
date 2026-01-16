package api

import (
	"net/http"
	"time"

	"github.com/Samu-Amy/Shokora/internal/store"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// - Structs -
type App struct {
	config Config
	router *chi.Mux
	store  *store.Storage
	logger *zap.SugaredLogger
}

type Config struct {
	Addr string
	Db   DbConfig
	Mail MailConfig
	// Env  string
}

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type MailConfig struct {
	EmailVerificationTokenExp time.Duration
	PasswordResetTokenExp     time.Duration
}

// - Functions/Methods -
func NewApp(config Config, store *store.Storage, logger *zap.SugaredLogger) *App {
	app := &App{
		config: config,
		store:  store,
		logger: logger,
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
