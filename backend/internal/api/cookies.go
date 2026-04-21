package api

import (
	"net/http"
	"time"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
)

const (
	cookiePath             string = "/api"
	AccessTokenCookieName  string = "access_token"
	RefreshTokenCookieName string = "refresh_token"
)

func (app *App) setAuthCookies(w http.ResponseWriter, authTokensDto payloads.AuthTokensDto) {

	accessCookie := newSecureCookie(AccessTokenCookieName, authTokensDto.AccessToken, authTokensDto.AccessTokenExpiresAt)
	refreshCookie := newSecureCookie(RefreshTokenCookieName, authTokensDto.PlainRefreshToken, authTokensDto.RefreshTokenExpiresAt)

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)
}

func (app *App) clearAuthCookies(w http.ResponseWriter) {
	accessCookie := expiredSecureCookie(AccessTokenCookieName)
	refreshCookie := expiredSecureCookie(RefreshTokenCookieName)

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)
}

// ----- UTILS -----

func newSecureCookie(name, value string, expiration time.Time) http.Cookie {
	return http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode, // TODO: forntend e backend devono avere stesso dominio (con reverse proxy di nginx dovrebbe andare bene)
		MaxAge:   getMaxAge(expiration),
		Expires:  expiration.UTC(),
		Path:     cookiePath,
	}
}

func expiredSecureCookie(name string) http.Cookie {
	return http.Cookie{
		Name:     name,
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		Path:     cookiePath,
	}
}

func getMaxAge(expiration time.Time) int {
	maxAge := int(time.Until(expiration).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}
	return maxAge
}
