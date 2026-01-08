package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// - Structs -
type App struct {
	config Config
	router *chi.Mux
}

type Config struct {
	Addr string
}

// - Functions/Methods -
func NewApp(config Config) *App {
	app := &App{
		config: config,
		router: initRouter(),
	}

	return app
}

func (app *App) Run() error {
	server := http.Server{
		Addr:    app.config.Addr,
		Handler: app.router,
	}

	fmt.Printf("Listening on localhost%s", app.config.Addr)
	return server.ListenAndServe()
}
