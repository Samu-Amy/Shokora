package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/store"
)

// - Return an error -
func (app *App) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *App) requestTimeoutError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("request timeout error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusRequestTimeout, "failed to process request in time")
}

func (app *App) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *App) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *App) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("conflict error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusConflict, "conflict")
}

// - Parse error -
func (app *App) parseError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		app.requestTimeoutError(w, r, err)

	case errors.Is(err, store.ErrNotFound):
		app.notFoundError(w, r, err)

	case errors.Is(err, store.ErrVersionConlflict):
		app.conflictError(w, r, err)

	default:
		app.internalServerError(w, r, err)
	}
}
