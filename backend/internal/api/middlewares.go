package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Samu-Amy/Shokora/internal/store"
	"github.com/golang-jwt/jwt/v5"
)

// - Rate Limiter -

func (app *App) rateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.RateLimiter.Enabled {
			if allow, retryAfter := app.rateLimiter.Allow(r.RemoteAddr); !allow {
				app.rateLimitExceededError(w, r, retryAfter.String())
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// - Authentication -

func (app *App) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		// Get Auth header
		authHeader := r.Header.Get("Authorization") // TODO: fai functions utils (?)
		if authHeader == "" {
			app.unauthorizedError(w, r, ErrTokenInvalid)
			return
		}

		// Parse Auth header ("authorization: Bearer <token>")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedError(w, r, ErrTokenInvalid)
			return
		}

		token := parts[1]

		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedError(w, r, ErrTokenExpired)
			return
		}

		// Get user id
		claims := jwtToken.Claims.(jwt.MapClaims)

		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		// Get user
		user, err := app.store.User.GetById(ctx, userId) // TODO: gestione caso utente non verificato (?)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		//* Save user in context
		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// - Authorization - Roles -

func (app *App) employeeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get User
		user, ok := getUserFromContext(r)
		if !ok || user == nil {
			app.unauthorizedError(w, r, errors.New("user not found"))
			return
		}

		// Check User Role
		if !user.IsRoleValid(store.RoleEmployee) {
			app.forbiddenError(w, r, errors.New("user doesn't have the permission"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *App) adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO: aggiungi controllo permessi (solo per employee)

		// Get User
		user, ok := getUserFromContext(r)
		if !ok || user == nil {
			app.unauthorizedError(w, r, errors.New("user not found"))
			return
		}

		// Check User Role
		if !user.IsRoleValid(store.RoleAdmin) {
			app.forbiddenError(w, r, errors.New("user doesn't have the permission"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *App) devMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get User
		user, ok := getUserFromContext(r)
		if !ok || user == nil {
			app.unauthorizedError(w, r, errors.New("user not found"))
			return
		}

		// Check User Role
		if !user.IsRoleValid(store.RoleDev) {
			app.forbiddenError(w, r, errors.New("user doesn't have the permission"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// - Authorization - Ownership -

func (app *App) userOwnershipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get user id
		userId, err := app.getIdFromParam(r, userIdParam)
		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		// Get User
		user, ok := getUserFromContext(r)
		if !ok || user == nil {
			app.unauthorizedError(w, r, errors.New("user not found"))
			return
		}

		// Check User Id
		if user.Id != userId {
			app.forbiddenError(w, r, errors.New("trying to update other user's data")) // TODO: cambia (?)
			return
		}

		next.ServeHTTP(w, r)
	})
}
