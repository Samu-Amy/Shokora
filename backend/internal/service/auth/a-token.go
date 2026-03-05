package authservice

import (
	"context"
	"strconv"
	"time"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	interrors "github.com/Samu-Amy/Shokora/internal/errors/int"
	"github.com/golang-jwt/jwt/v5"
)

/*
Check the access token and return an *AuthTokensCheckDto if valid

Returs error:
  - if error is interrors.IErrUnauthorized -> ! MUST RETURN early from auth check (the two tokens doesn't correspond -> something is wrong) !
  - else is ok to coninue checking the refresh token
*/
func (service *AuthService) checkAccessToken(ctx context.Context, accessToken string, hashedRefreshToken []byte) (*payloads.AuthTokensCheckDto, error) {
	// Validate and obtain jwt
	jwtToken, err := service.jwtAuthenticator.ValidateJWTToken(accessToken)
	if err != nil || jwtToken == nil {
		service.logger.Info("JWT not valid", "error", err)
		return nil, interrors.IErrInvalid
	}

	// Get user Id from claims
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		service.logger.Info("Invalid jwt claims type", "error", err)
		return nil, interrors.IErrInvalid
	}

	subject, err := claims.GetSubject()
	if err != nil {
		service.logger.Info("Error getting subject from jwt", "error", err)
		return nil, interrors.IErrInvalid
	}

	// Get userId from jwt
	userId, err := strconv.ParseInt(subject, 10, 64)
	if err != nil {
		service.logger.Info("Error parsing userId from jwt subject", "error", err)
		return nil, interrors.IErrInvalid
	}

	// Get sessionId and userId from db (refresh token) //? (maggiore sicurezza, performance minori, l'alternativa è mettere sessionId nel jwt dell'access token (per evitare query db))
	sessionData, err := service.refreshTokenRepo.GetSessionDataByToken(ctx, hashedRefreshToken)
	if err != nil {
		service.logger.Info("Error getting the session id for the refresh token", "error", err)
		return nil, interrors.IErrInvalid
	}

	// Check userId coherence between Access Token (jwt) and Refresh Token (db)
	if userId != sessionData.UserId {
		service.logger.Warnw("UserIds from access and refresh token doesn't correspond")

		// Delete session
		err = service.userSessionRepo.Delete(ctx, sessionData.SessionId)
		if err != nil {
			service.logger.Warnw("Error deleting Session", "error", err)
		}

		return nil, interrors.IErrUnauthorized
	}

	// Create payload with data
	authTokensCheckDto := payloads.AuthTokensCheckDto{
		IsAccessTokenValid: true,
		UserId:             sessionData.UserId,
		SessionId:          sessionData.SessionId,
	}

	return &authTokensCheckDto, nil
}

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
