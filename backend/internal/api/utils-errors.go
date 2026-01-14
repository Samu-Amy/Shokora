package api

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/store/postgres"
)

// - Return an error -
func (app *App) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *App) requestTimeoutError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("request timeout error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusRequestTimeout, "failed to process request in time")
}

func (app *App) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *App) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *App) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("conflict error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusConflict, "conflict")
}

// - Parse error -
func (app *App) parseError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		app.requestTimeoutError(w, r, err)

	case errors.Is(err, postgres.ErrNotFound):
		app.notFoundError(w, r, err)

	case errors.Is(err, postgres.ErrVersionConlflict):
		app.conflictError(w, r, err)

	default:
		app.internalServerError(w, r, err)
	}
}
