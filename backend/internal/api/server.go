package api

import (
	"log"
	"net/http"
	"time"

	"github.com/Samu-Amy/Shokora/internal/store"
	"github.com/go-chi/chi/v5"
)

// - Structs -
type App struct {
	config Config
	router *chi.Mux
	store  *store.Storage
}

type Config struct {
	Addr string
	Db   DbConfig
}

type DbConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

// - Functions/Methods -
func NewApp(config Config, store *store.Storage) *App {
	app := &App{
		config: config,
		store:  store,
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

	log.Printf("Server started, listening on localhost%s", app.config.Addr)
	return server.ListenAndServe()
}
