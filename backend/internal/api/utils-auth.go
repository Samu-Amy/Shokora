package api

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (app *App) setAuthCookies(w http.ResponseWriter, userId int64) error {

	timeNow := time.Now()
	accessTokenExp := timeNow.Add(app.config.Auth.Token.AccessTokenExp)

	// Generate Access Token (and add claims)
	claims := jwt.MapClaims{
		"sub": userId, // subject
		"exp": accessTokenExp.Unix(),
		"iat": timeNow.Unix(),                 // issued at
		"nbf": timeNow.Unix(),                 // not before time
		"iss": app.config.Auth.Token.Issuer,   // issuer
		"aud": app.config.Auth.Token.Audience, // audience
	}

	accessToken, err := app.jwtAuthenticator.GenerateJWTToken(claims)
	if err != nil {
		return err
	}

	// TODO: crea refresh token
	refreshToken := ""

	// Create and set cookie
	accessCookie := http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(app.config.Auth.Token.AccessTokenExp.Seconds()),
		Path:     "/api",
	}

	refreshCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(app.config.Auth.Token.RefreshTokenExp.Seconds()), // TODO: aggiornare quando si aumenta la durata (usa l'expiry settata sul db)
		Path:     "/api",
	}

	http.SetCookie(w, &accessCookie)
	http.SetCookie(w, &refreshCookie)

	return nil
}
