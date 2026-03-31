package authservice

import (
	"context"

	"github.com/Samu-Amy/Shokora/internal/auth"
	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
)

func (service *AuthService) GenerateOAuthUrl(ctx context.Context) (string, error) {

	// Generate state
	state, err := auth.GenerateBase64Token(32)
	if err != nil {
		return "", domerrors.ErrInternalError
	}

	// Save state in db
	err = service.oAuthStateRepo.Create(ctx, state)
	if err != nil {
		return "", domerrors.ParseIntError(err)
	}

	// Generate url
	oAuthUrl := service.config.Auth.GoogleOAuthConfig.AuthCodeURL(state)

	return oAuthUrl, nil
}
