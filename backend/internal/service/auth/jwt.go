package authservice

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: implementa

func (service *AuthService) generateAccessToken(userId int64) (*string, error) {

	// Set times
	timeNow := time.Now()
	accessTokenExp := time.Now().Add(service.config.Token.AccessTokenExp)

	// Generate Access Token (and add claims)
	claims := jwt.MapClaims{
		"sub": userId, // subject
		"exp": accessTokenExp.Unix(),
		"iat": timeNow.Unix(),                // issued at
		"nbf": timeNow.Unix(),                // not before time
		"iss": service.config.Token.Issuer,   // issuer
		"aud": service.config.Token.Audience, // audience
	}

	accessToken, err := service.jwtAuthenticator.GenerateJWTToken(claims)
	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}
