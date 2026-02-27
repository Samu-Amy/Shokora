package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
	user_repo "github.com/Samu-Amy/Shokora/internal/store/user"
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
		authHeader := r.Header.Get(authHeader)
		if authHeader == "" {
			app.unauthorizedError(w, r, domerrors.ErrInvalid)
			return
		}

		// Parse Auth header ("authorization: Bearer <token>")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != bearer {
			app.unauthorizedError(w, r, domerrors.ErrInvalid)
			return
		}

		token := parts[1]

		jwtToken, err := app.jwtAuthenticator.ValidateJWTToken(token) // TODO: usa service (gestione sia di Access che di Refresh tokens)
		if err != nil {
			app.unauthorizedError(w, r, domerrors.ErrUnauthorized)
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
		user, err := app.service.User.GetById(ctx, userId)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		// Check if user is not blocked
		if !user.IsActive {
			app.unauthorizedError(w, r, err) // TODO: usa errore dedicato (bloccato)
			return
		}

		//* Save user in context
		ctxWithUser := context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctxWithUser))
	})
}

// - Authorization - Ownership - (non serve più)

// // Verify that the data the user is trying to access is theirs.
// // Must be called after the authMiddleware
// func (app *App) userDataOwnershipMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		// Get user id from parameters
// 		userId, err := app.getIdFromParam(r, userIdParam)
// 		if err != nil {
// 			app.badRequestError(w, r, err)
// 			return
// 		}

// 		// Get User from context (auth middleware)
// 		user, ok := getUserFromContext(r)
// 		if !ok || user == nil {
// 			app.unauthorizedError(w, r, ErrUserNotFound)
// 			return
// 		}

// 		// Check User Id
// 		if user.Id != userId {
// 			app.forbiddenError(w, r, ErrUserNotAuthorized) // TODO: cambia (?)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// - Authorization - User Verified -

// Verify that the user's email is verified.
// Must be called after the authMiddleware
func (app *App) userVerifiedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get User from context (auth middleware)
		user, ok := getUserFromContext(r)
		if !ok || user == nil {
			app.unauthorizedError(w, r, domerrors.ErrUnauthorized)
			return
		}

		// Check if User is verified
		if !user.IsVerified {
			app.unauthorizedError(w, r, domerrors.ErrUnauthorized) // TODO: in forntend chiedi verifica
			return
		}

		next.ServeHTTP(w, r)
	})
}

// - Authorization - Roles -

// Verify that the user's role is >= employee.
// Must be called after the authMiddleware
func (app *App) employeeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get User
		user, ok := getUserFromContext(r)
		if !ok || user == nil {
			app.unauthorizedError(w, r, domerrors.ErrUnauthorized)
			return
		}

		// Check User Role
		if !user.IsRoleValid(user_repo.RoleEmployee) {
			app.forbiddenError(w, r, domerrors.ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Verify that the user's role is >= admin.
// Must be called after the authMiddleware
func (app *App) adminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO: aggiungi controllo permessi (solo per employee)

		// Get User
		user, ok := getUserFromContext(r)
		if !ok || user == nil {
			app.unauthorizedError(w, r, domerrors.ErrUnauthorized)
			return
		}

		// Check User Role
		if !user.IsRoleValid(user_repo.RoleAdmin) {
			app.forbiddenError(w, r, domerrors.ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Verify that the user's role is >= dev.
// Must be called after the authMiddleware
func (app *App) devMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get User
		user, ok := getUserFromContext(r)
		if !ok || user == nil {
			app.unauthorizedError(w, r, domerrors.ErrUnauthorized)
			return
		}

		// Check User Role
		if !user.IsRoleValid(user_repo.RoleDev) {
			app.forbiddenError(w, r, domerrors.ErrForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
