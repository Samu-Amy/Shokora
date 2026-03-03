package authservice

import (
	"time"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/golang-jwt/jwt/v5"
)

// Generate an Access Token (JWT) with expiration and add them to authTokensTdo
func (service *AuthService) addJWTAccessToken(authTokensTdo *payloads.AuthTokensDto, userId int64) error {

	// Set expiration
	timeNow := time.Now()
	accessTokenExpiresAt := time.Now().Add(service.config.Token.AccessTokenExp)

	// Create Claims
	claims := jwt.MapClaims{
		"sub": userId, // subject
		"exp": accessTokenExpiresAt.Unix(),
		"iat": timeNow.Unix(),                // issued at
		"nbf": timeNow.Unix(),                // not before time
		"iss": service.config.Token.Issuer,   // issuer
		"aud": service.config.Token.Audience, // audience
	}

	// Generate Access Token (and add claims)
	accessToken, err := service.jwtAuthenticator.GenerateJWTToken(claims)
	if err != nil {
		service.logger.Warnw("Error generating access token (jwt)", "error", err)
		return err
	}

	// Add token and expiration to authTokenDto
	authTokensTdo.AccessToken = accessToken
	authTokensTdo.AccessTokenExpiresAt = accessTokenExpiresAt

	return nil
}
