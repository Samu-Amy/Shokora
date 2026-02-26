package api

import (
	"net/http"
	"strconv"

	"github.com/Samu-Amy/Shokora/internal/store/user"
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
func getUserFromContext(r *http.Request) (*user.User, bool) {
	user, ok := r.Context().Value(userCtx).(*user.User)
	return user, ok
}

// - Params -

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

// - Auth -
func (app *App) hashPassword(plainPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), app.config.Auth.PasswordHashingCost)
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
