package api

import (
	"context"
	"net/http"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

func (app *App) setAuthCookies(w http.ResponseWriter, userId int64, plainRefreshToken string, refreshTokenExp time.Time) error {

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

	// Create and set cookies
	accessCookie := http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(time.Until(accessTokenExp).Seconds()),
		Path:     "/api",
	}

	refreshMaxAge := int(time.Until(refreshTokenExp).Seconds())
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

func (app *App) generateRefreshToken(ctx context.Context, userId int64) (*auth.RefreshTokenPayload, error) {
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
		SessionId:   sessionId,
		HashedToken: hashedToken,
		Exp:         app.config.Auth.Token.RefreshTokenExp,
	}

	// Save token in
	tokenExpiresAt, err := app.store.RefreshTokens.CreateToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return &auth.RefreshTokenPayload{
		PlainToken: *token,
		ExpiresAt:  tokenExpiresAt,
	}, nil
}
