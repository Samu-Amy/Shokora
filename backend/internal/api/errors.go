package api

import (
	"context"
	"errors"
	"net/http"

	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
)

// TODO: qua devono esserci solo domerrors (il service deve gestire gli interrors e tradurli in domerrors) - evita comunque di inviare interrors se dovessero esserci per sbaglio (controlla manualmente gli errori com'è già adesso)

// ----- SEND ERROR -----

// - Internal Errors (fixed message) -

func (app *App) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	}

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *App) rateLimitExceededError(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("rate limit exceeded (too many requests) error", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	writeJSONError(w, http.StatusTooManyRequests, "too many requests")
}

func (app *App) requestTimeoutError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		app.logger.Warnf("request timeout error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	}

	writeJSONError(w, http.StatusRequestTimeout, "failed to process request in time")
}

func (app *App) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		app.logger.Errorf("conflict error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	}

	writeJSONError(w, http.StatusConflict, "conflict")
}

// - Dynamic Errors (message from error) -

func (app *App) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		app.logger.Warnf("bad request error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	}

	errorMessage := "bad_request"

	if domerrors.IsDomainErr(err) {
		errorMessage = err.Error()
	}

	writeJSONError(w, http.StatusBadRequest, errorMessage)
}

func (app *App) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		app.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	}

	writeJSONError(w, http.StatusNotFound, "not_found")
}

func (app *App) unauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		app.logger.Warnf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	}

	errorMessage := "unauthorized"

	if domerrors.IsDomainErr(err) {
		errorMessage = err.Error()
	}

	writeJSONError(w, http.StatusUnauthorized, errorMessage)
}

func (app *App) forbiddenError(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		app.logger.Warnf("forbidden error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	}

	errorMessage := "forbidden"

	if domerrors.IsDomainErr(err) {
		errorMessage = err.Error()
	}

	writeJSONError(w, http.StatusForbidden, errorMessage)
}

// ----- PARSE ERROR -----

// TODO: aggiorna con tutti i domerrors (e ricontrolla assegnazioni)
func (app *App) parseError(w http.ResponseWriter, r *http.Request, err error) {
	switch {

	// - Internal -
	case !domerrors.IsDomainErr(err):
		app.internalServerError(w, r, err)

	case errors.Is(err, domerrors.ErrMaxRetriesExceeded): // Better to handle this case by case
		app.internalServerError(w, r, err)

	// - Timeout -
	case errors.Is(err, context.DeadlineExceeded):
		app.requestTimeoutError(w, r, err)

	// - Not Found -
	case errors.Is(err, domerrors.ErrNotFound):
		app.notFoundError(w, r, err)

	// - Conflict -
	case errors.Is(err, domerrors.ErrConflict):
		app.conflictError(w, r, err)

	// - Bad Request -

	// Password
	case errors.Is(err, domerrors.ErrCommonPassword), errors.Is(err, domerrors.ErrSamePassword), errors.Is(err, domerrors.ErrDifferentPassword),
		errors.Is(err, domerrors.ErrDuplicateEmail), errors.Is(err, domerrors.ErrInvalid), errors.Is(err, domerrors.ErrInvalidName), errors.Is(err, domerrors.ErrInvalidDate):
		app.badRequestError(w, r, err)

	// - Unauthorized -
	case errors.Is(err, domerrors.ErrUnauthorized), errors.Is(err, domerrors.ErrNotVerified):
		app.unauthorizedError(w, r, err)

	// - Forbidden -
	case errors.Is(err, domerrors.ErrMaxAttemptsExceeded):
		app.forbiddenError(w, r, err)

	default:
		app.internalServerError(w, r, err)
	}
}
