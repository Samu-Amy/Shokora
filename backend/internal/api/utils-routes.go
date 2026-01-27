package api

import (
	"net/http"
	"strconv"

	"github.com/Samu-Amy/Shokora/internal/store"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

// - Context -

// Keys
type contextKey uint8

const (
	userCtx contextKey = iota
)

// Functions
func getUserFromContext(r *http.Request) (*store.User, bool) {
	user, ok := r.Context().Value(userCtx).(*store.User)
	return user, ok
}

// - Params -

// Constants
const userIdParam = "userId"
const productIdParam = "productId"

// Methods
func (app *App) getIdFromParam(r *http.Request, idParamName string) (int64, error) {
	idParam := chi.URLParam(r, idParamName)

	resourceId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return -1, err
	}

	return resourceId, nil
}

// - Auth -
func (app *App) hashPassword(plainPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), app.config.Auth.HashingCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

// func (app *App) setAuthCookie(w http.ResponseWriter, token string) {
// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "auth_token",
// 		Value:    token,
// 		Path:     "/",
// 		MaxAge:   int(app.config.Auth.Token.Exp.Seconds()),
// 		HttpOnly: true,
// 		Secure:   app.config.Env == "production", // true in prod
// 		SameSite: http.SameSiteStrictMode,
// 	})
// }
