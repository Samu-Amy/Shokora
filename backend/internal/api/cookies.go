package api

import (
	"net/http"
	"time"
)

func (app *App) setAuthCookies(w http.ResponseWriter, userId int64, accessToken, plainRefreshToken string, accessTokenExpiresAt, refreshTokenExpiresAt time.Time) error {

	// Create and set cookies
	accessCookie := http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(time.Until(refreshTokenExpiresAt).Seconds()),
		Path:     "/api",
	}

	refreshMaxAge := int(time.Until(refreshTokenExpiresAt).Seconds())
	if refreshMaxAge < 0 {
		refreshMaxAge = 0
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    plainRefreshToken, // Plain token
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   refreshMaxAge,
		Path:     "/api",
	}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)

	return nil
}
