package authservice

import (
	"context"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
)

/*
Create Session, Refresh Token and Access Token, return the AuthTokensDto for setting the cookies
*/
func (service *AuthService) createAuthTokens(ctx context.Context, userId int64) (*payloads.AuthTokensDto, error) {

	// Create Refresh Token
	authTokensDto, sessionId, err := service.createNewSessionAndRefreshToken(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Create Access Token
	err = service.addJWTAccessToken(authTokensDto, sessionId, userId)
	if err != nil {
		return nil, err
	}

	return authTokensDto, nil
}
