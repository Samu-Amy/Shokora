package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/store"
)

var (
	ErrUserBlocked  = errors.New("user_blocked")
	ErrTokenInvalid = errors.New("token_invalid")
	ErrTokenExpired = errors.New("token_expired")
	// TODO: aggiungere versioni per refresh token?
)

// TODO: passa errori strutturati (sopra) al frontend (invece che hardcoded strings)

// - Return an error -
func (app *App) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *App) rateLimitExceededError(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("rate limit exceeded (too many requests) error", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	writeJSONError(w, http.StatusTooManyRequests, "too many requests")
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

func (app *App) unauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *App) forbiddenError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("forbiddeb error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusForbidden, "forbidden")
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

	case errors.Is(err, store.ErrDuplicateEmail):
		app.badRequestError(w, r, err) // TODO: passare quale dato è duplicato (email) per poter mostrare un messaggio più preciso all'utente (?)

	case errors.Is(err, store.ErrUnauthorized):
		app.unauthorizedError(w, r, err)

	default:
		app.internalServerError(w, r, err)
	}
}
