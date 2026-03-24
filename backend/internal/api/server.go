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
	"github.com/Samu-Amy/Shokora/internal/config"
	"github.com/Samu-Amy/Shokora/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// - Structs -
type App struct {
	config      config.Config
	router      *chi.Mux
	service     *service.Service
	logger      *zap.SugaredLogger
	rateLimiter ratelimiter.RateLimiterI
}

// - Functions/Methods -
func NewApp(
	config config.Config,
	service *service.Service,
	logger *zap.SugaredLogger,
	rateLimiter ratelimiter.RateLimiterI,
) *App {
	app := &App{
		config:      config,
		service:     service,
		logger:      logger,
		rateLimiter: rateLimiter,
	}
	app.router = app.InitRouter()

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

// func NewMockApp(store *store.Storage, logger *zap.SugaredLogger, authenticator auth.Authenticator) *App {
// 	app := &App{
// 		store:         store,
// 		logger:        logger,
// 		authenticator: authenticator,
// 	}
// 	app.router = app.initRouter()

// 	return app
// }

// func (app *App) GetRouter() *chi.Mux {
// 	return app.router
// }

// func (app *App) GetAuthenticator() auth.Authenticator {
// 	return app.authenticator
// }
