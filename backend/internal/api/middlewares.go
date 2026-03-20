package api

import (
	"context"
	"net/http"

	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
	user_repo "github.com/Samu-Amy/Shokora/internal/store/user"
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

		// Get Access Token
		var accessToken string
		accessCookie, err := r.Cookie(accessTokenCookieName)
		if err == nil {
			accessToken = accessCookie.Value // The Access Token can be expired (and the cookie deleted), the important one is the refresh
		}

		// Get Refresh Token
		refreshCookie, err := r.Cookie(refreshTokenCookieName)
		if err != nil {
			app.unauthorizedError(w, r, domerrors.ErrUnauthorized)
			return
		}

		plainRefreshToken := refreshCookie.Value

		// Check Tokens
		authTokensCheckDto, isAccessTokenValid, err := app.service.Auth.HandleAuthTokensCheck(ctx, accessToken, plainRefreshToken)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		// Set cookies (if tokens are updated)
		if !isAccessTokenValid {
			app.setAuthCookies(w, authTokensCheckDto.AuthTokensDto)
		}

		// Get user
		user, err := app.service.User.GetById(ctx, authTokensCheckDto.UserId)
		if err != nil {
			app.unauthorizedError(w, r, err)
			return
		}

		// Check if user is not blocked
		if !user.IsActive {
			app.forbiddenError(w, r, domerrors.ErrForbidden)
			return
		}

		//* Save user and sessionId in context
		ctxWithUser := context.WithValue(ctx, userCtx, user)
		ctxWithSessionId := context.WithValue(ctxWithUser, sessionIdCtx, authTokensCheckDto.SessionId)

		next.ServeHTTP(w, r.WithContext(ctxWithSessionId))
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

/*
Verify that the user's email is verified.
Must be called after the authMiddleware
*/
func (app *App) userVerifiedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get User from context (auth middleware)
		user, ok := r.Context().Value(userCtx).(*user_repo.User)
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

/*
Verify that the user's role is >= requiredRole.
Must be called after the authMiddleware
*/
func (app *App) roleMiddleware(requiredRole user_repo.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Get User
			user, ok := r.Context().Value(userCtx).(*user_repo.User)
			if !ok || user == nil {
				app.unauthorizedError(w, r, domerrors.ErrUnauthorized)
				return
			}

			// Check if User is verified
			if !user.IsVerified {
				app.unauthorizedError(w, r, domerrors.ErrUnauthorized) // TODO: in forntend chiedi verifica
				return
			}

			// Check User Role
			if !user.IsRoleValid(requiredRole) {
				app.forbiddenError(w, r, domerrors.ErrForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// - Authorization - Permissions -

/*
Verify that the user's role is > requiredRole or it is == to requiredRole but with required permission.
Must be called after the authMiddleware
*/
func (app *App) permissionMiddleware(requiredRole user_repo.Role, requiredPermission user_repo.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Get User
			user, ok := r.Context().Value(userCtx).(*user_repo.User)
			if !ok || user == nil {
				app.unauthorizedError(w, r, domerrors.ErrUnauthorized)
				return
			}

			// Check if User is verified
			if !user.IsVerified {
				app.unauthorizedError(w, r, domerrors.ErrUnauthorized) // TODO: in forntend chiedi verifica
				return
			}

			// Check User Role
			if !user.IsRoleValid(requiredRole) {
				app.forbiddenError(w, r, domerrors.ErrForbidden)
				return
			}

			// Check permission
			if user.Role == (requiredRole) {
				if !user.HasPermission(requiredPermission) {
					app.forbiddenError(w, r, domerrors.ErrForbidden)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
