package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/errorcodes"
)

// TODO: passa errori strutturati (sopra) al frontend (invece che hardcoded strings)

// - Internal Errors (fixed message) -
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

func (app *App) unauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

// - Dynamic Errors (message from error) -

func (app *App) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	errorMessage := "bad_request"

	// Frontend receives "ErrInvalid"
	if errors.Is(err, errorcodes.InternalErrExpired) {
		err = errorcodes.ErrInvalid
	}

	if errors.Is(err, errorcodes.ErrDuplicateEmail) || errors.Is(err, errorcodes.ErrInvalid) {
		errorMessage = err.Error()
	}

	writeJSONError(w, http.StatusBadRequest, errorMessage)
}

func (app *App) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusNotFound, "not_found")
}

func (app *App) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("conflict error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJSONError(w, http.StatusConflict, "conflict")
}

func (app *App) forbiddenError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("forbidden error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	errorMessage := "forbidden"

	if errors.Is(err, errorcodes.ErrMaxAttemptsExceeded) {
		errorMessage = err.Error()
	}

	writeJSONError(w, http.StatusForbidden, errorMessage)
}

// TODO: attenzione agli errori che passano al frontend (bisogna usare solo quelli di errorcodes per evitare leaks di informazioni)

// - Parse error -
func (app *App) parseError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		app.requestTimeoutError(w, r, err)

	case errors.Is(err, errorcodes.ErrNotFound):
		app.notFoundError(w, r, err)

	case errors.Is(err, errorcodes.ErrConflict):
		app.conflictError(w, r, err)

	case errors.Is(err, errorcodes.ErrDuplicateEmail), errors.Is(err, errorcodes.ErrInvalid), errors.Is(err, errorcodes.InternalErrExpired):
		app.badRequestError(w, r, err)

	case errors.Is(err, errorcodes.ErrUnauthorized), errors.Is(err, errorcodes.ErrNotVerified):
		app.unauthorizedError(w, r, err)

	case errors.Is(err, errorcodes.ErrMaxAttemptsExceeded):
		app.forbiddenError(w, r, err)

	// Better to handle these case by case
	case errors.Is(err, errorcodes.ErrMaxRetriesExceeded), errors.Is(err, errorcodes.ErrVerification), errors.Is(err, errorcodes.ErrEmailNotSent):
		app.internalServerError(w, r, err)

	default:
		app.internalServerError(w, r, err)
	}
}
