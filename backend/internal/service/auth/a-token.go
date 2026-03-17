package authservice

import (
	"errors"
	"time"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/auth"
	interrors "github.com/Samu-Amy/Shokora/internal/errors/int"
	"github.com/golang-jwt/jwt/v5"
)

/*
Check the access token and return an *AuthTokensCheckDto if valid

Returs error:
  - if error is interrors.IErrUnauthorized -> ! MUST RETURN early from auth check (the two tokens doesn't correspond -> something is wrong) !
  - else is ok to coninue checking the refresh token
*/
func (service *AuthService) checkAccessToken(accessToken string) (*payloads.AuthTokensCheckDto, error) {

	// Validate and obtain claims
	claims, err := service.jwtAuthenticator.ValidateJWTToken(accessToken)
	if err != nil {
		service.logger.Infow("JWT not valid", "error", err)

		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, interrors.IErrExpired
		}

		return nil, interrors.IErrInvalid
	}

	// Get userId from claims
	userId := claims.UserId
	if userId <= 0 {
		service.logger.Info("UserId from jwt claims cannot be <= 0")
		return nil, interrors.IErrInvalid
	}

	// Get sessionId from claims
	sessionId := claims.SessionId
	if sessionId <= 0 {
		service.logger.Info("SessionId from jwt claims cannot be <= 0")
		return nil, interrors.IErrInvalid
	}

	// Create payload with data
	authTokensCheckDto := payloads.AuthTokensCheckDto{
		IsAccessTokenValid: true,
		UserId:             userId,
		SessionId:          sessionId,
	}

	return &authTokensCheckDto, nil
}

// Generate an Access Token (JWT) with expiration and add them to authTokensTdo
func (service *AuthService) addJWTAccessToken(authTokensDto *payloads.AuthTokensDto, sessionId int64, userId int64) error {

	// Set expiration
	timeNow := time.Now().UTC()
	accessTokenExpiresAt := timeNow.Add(service.config.Token.AccessTokenExp)

	// Create Claims
	claims := auth.UserClaims{
		UserId:    userId,
		SessionId: sessionId,
		RegisteredClaims: jwt.RegisteredClaims{
			// Subject:   strconv.FormatInt(userId, 10),
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiresAt),
			IssuedAt:  jwt.NewNumericDate(timeNow),
			NotBefore: jwt.NewNumericDate(timeNow),
			Issuer:    service.config.Token.Issuer,
			Audience:  []string{service.config.Token.Audience},
		},
	}

	// Generate Access Token (and add claims)
	accessToken, err := service.jwtAuthenticator.GenerateJWTToken(claims)
	if err != nil {
		service.logger.Warnw("Error generating access token (jwt)", "error", err)
		return err
	}

	// Add token and expiration to authTokenDto
	authTokensDto.AccessToken = accessToken
	authTokensDto.AccessTokenExpiresAt = accessTokenExpiresAt

	return nil
}
