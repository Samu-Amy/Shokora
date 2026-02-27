package api

import (
	"net/http"
	"time"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
)

func (app *App) setAuthCookies(w http.ResponseWriter, authTokensDto payloads.AuthTokensDto) {

	// Create and set cookies
	accessCookie := newSecureCookie("access_token", authTokensDto.AccessToken, authTokensDto.AccessTokenExpiresAt)

	refreshCookie := newSecureCookie("refresh_token", authTokensDto.PlainRefreshToken, authTokensDto.RefreshTokenExpiresAt)

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)
}

func newSecureCookie(name, value string, expiration time.Time) http.Cookie {
	return http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode, // TODO: forntend e backend devono avere stesso dominio (con reverse proxy di nginx dovrebbe andare bene)
		MaxAge:   getMaxAge(expiration),
		Expires:  expiration.UTC(),
		Path:     "/api",
	}
}

func getMaxAge(expiration time.Time) int {
	maxAge := int(time.Until(expiration).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}
	return maxAge
}
