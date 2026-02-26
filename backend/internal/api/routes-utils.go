package api

import (
	"net/http"
	"strconv"

	"github.com/Samu-Amy/Shokora/internal/store/user"
	"github.com/go-chi/chi/v5"
)

// ----- HEADERS -----
const (
	AUTH_HEADER string = "Authorization"
	BEARER      string = "Bearer"
)

// ----- CONTEXT -----

// Keys
type contextKey uint8

const (
	userCtx contextKey = iota
)

// Functions
func getUserFromContext(r *http.Request) (*user.User, bool) {
	user, ok := r.Context().Value(userCtx).(*user.User)
	return user, ok
}

// ----- PARAMS -----

// Constants
const userIdParam = "userId"
const productIdParam = "productId"
const verificationTokenParam = "token"

// Methods
func (app *App) getInt64FromParam(r *http.Request, idParamName string) (int64, error) {
	param := chi.URLParam(r, idParamName)

	resourceId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return -1, err
	}

	return resourceId, nil
}
