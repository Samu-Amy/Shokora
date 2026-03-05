package authservice

import (
	"context"

	session "github.com/Samu-Amy/Shokora/internal/store/user-session"
)

/*
Get session and user Id of the session associated with the refresh token
*/
func (service *AuthService) getSessionData(ctx context.Context, hashedRefreshToken []byte) (*session.SessionData, bool) {
	sessionData, err := service.refreshTokenRepo.GetSessionDataByToken(ctx, hashedRefreshToken)
	if err != nil {
		service.logger.Warnw("Error getting the session id for the refresh token", "error", err)
		return nil, false
	}

	return sessionData, true
}
