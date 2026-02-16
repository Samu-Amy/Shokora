package api

import (
	"net/http"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
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

	// TODO: va bene solo per creazione nuova sessione (register e login/2fa) ma nel caso di refresh?
	// TODO: se login fare pulizia (eliminare token di sessioni scadute - attenzione agli expires aggiornati (vecchi token scaduti ma nuovi no -> sessione ancora valida), controlla per tutta la sessione)?
	refreshToken, err := app.generateRefreshToken(userId)
	if err != nil {
		return err
	}

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
		Value:    refreshToken, // Plain token
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

func (app *App) generateRefreshToken(userId int64) (*string, error) {
	token, err := auth.GenerateToken(app.config.Auth.Token.RefreshTokenByteSize)
	if err != nil {
		return nil, err
	}

	// Hash token and create Session Id
	hashedToken := auth.HashToken(token)

	// TODO: crea uuid session
	sessionId, err := auth.GenerateSessionId()
	if err != nil {
		return nil, err
	}

	refreshToken := auth.RefreshToken{
		UserId:      userId,
		SessionId:   *sessionId,
		HashedToken: hashedToken,
		Exp:         app.config.Auth.Token.RefreshTokenExp,
	}

	// Save token in
	app.store.

	return token, nil
}
