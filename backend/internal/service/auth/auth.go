package authservice

import (
	"context"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
)

/*
Create Session, Refresh Token and Access Token, return the AuthTokensDto for setting the cookies
*/
func (service *AuthService) createNewAuthTokens(ctx context.Context, userId int64) (*payloads.AuthTokensDto, error) {

	// Create Refresh Token
	authTokensCheckDto, err := service.createNewSessionAndRefreshToken(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Create Access Token
	err = service.addJWTAccessToken(authTokensCheckDto)
	if err != nil {
		return nil, err
	}

	return &authTokensCheckDto.AuthTokensDto, nil
}
