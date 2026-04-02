package authservice

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/auth"
	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
	interrors "github.com/Samu-Amy/Shokora/internal/errors/int"
	user_repo "github.com/Samu-Amy/Shokora/internal/store/user"
	"google.golang.org/api/idtoken"
)

// Generate URL
func (service *AuthService) GenerateGoogleOAuthUrl(ctx context.Context) (string, error) {

	// Generate state
	state, err := auth.GenerateBase64Token(32)
	if err != nil {
		service.logger.Warnw("Error generating state for Google OAuth ", "error", err)
		return "", domerrors.ErrInternalError
	}

	// Save state in db
	err = service.oAuthStateRepo.Create(ctx, state)
	if err != nil {
		service.logger.Warnw("Error creating state for Google OAuth in db ", "error", err)
		return "", domerrors.ParseIntError(err)
	}

	// Generate url
	oAuthUrl := service.config.Auth.GoogleOAuthConfig.AuthCodeURL(state)

	return oAuthUrl, nil
}

// Check and Authenticate user
func (service *AuthService) LoginUserWithGoogleOAuth(ctx context.Context, payload payloads.GoogleOAuthCallbackReq) (*payloads.LoginUserRes, *payloads.AuthTokensDto, error) {

	var loginUserRes payloads.LoginUserRes
	var authTokensDto payloads.AuthTokensDto
	var err error

	err = service.txManager.WithTx(ctx, func(tx *sql.Tx) error {

		// Check and delete state
		err = service.oAuthStateRepo.Delete(ctx, tx, payload.State)
		if err != nil {
			service.logger.Warnw("Error deleting state for Google OAuth from db ", "error", err)
			return err
		}

		// Get token from code
		token, err := service.config.Auth.GoogleOAuthConfig.Exchange(ctx, payload.Code)
		if err != nil {
			return domerrors.ErrInvalid
		}

		// Get data from token
		rawIDToken, ok := token.Extra("id_token").(string)
		if !ok {
			return domerrors.ErrInternalError
		}

		tokenPayload, err := idtoken.Validate(ctx, rawIDToken, service.config.Auth.GoogleOAuthConfig.ClientID)
		if err != nil {
			return domerrors.ErrInvalid
		}

		googleId := tokenPayload.Subject // "117382991234567890123"
		email, emailOk := tokenPayload.Claims["email"].(string)
		// name, nameOk := tokenPayload.Claims["name"].(string)
		firstName, firstNameOk := tokenPayload.Claims["given_name"].(string)
		lastName, lastNameOk := tokenPayload.Claims["family_name"].(string)
		// picture, _ := tokenPayload.Claims["picture"].(string)
		emailVerified, emailVerifiedOk := tokenPayload.Claims["email_verified"].(bool)

		// Get user by googleId (if exists, login)
		user, err := service.userRepo.GetByGoogleId(ctx, googleId) // TODO: serve FOR UPDATE ?
		if err != nil {
			// db errors
			if !errors.Is(err, interrors.IErrNotFound) {
				return err
			}

			if !emailOk || !emailVerifiedOk {
				return domerrors.ErrInvalid
			}

			// User not found by googleId -> Get user by email and check is verified (if exists and is verified, update google_id -> connect accounts)
			user, err = service.userRepo.GetByEmailForUpdate(ctx, tx, email)
			if err != nil {
				// db errors
				if !errors.Is(err, interrors.IErrNotFound) {
					return err
				}

				if !firstNameOk {
					return domerrors.ErrInvalid
				}

				// User does not exists -> Create user
				if !lastNameOk {
					lastName = ""
				}

				user = &user_repo.User{
					GoogleId:   googleId,
					FirstName:  firstName,
					LastName:   lastName,
					Email:      email,
					IsVerified: emailVerified,
				}

				err = service.userRepo.Create(ctx, tx, user)
				if err != nil {
					return err
				}
			}

			// Check if verified (if not -> error)
			if !user.IsVerified {
				return domerrors.ErrNotVerified // TODO: serve avere account verificato se si vuole accedere ad un account email e password usando google (ricorda di dirlo nel frontend)
			}

			// Add google_id (connect email and password and google accounts)
			err = service.userRepo.SetGoogleId(ctx, tx, user.Id, googleId) // soft error
			if err != nil {

			}
		}

		loginUserRes.User = payloads.ToUserRes(*user)

		// Create 2FA verification if needed
		// TODO: controlla se è richiesta 2fa ed aggiorna loginUserRes

		// Create authTokensDto

		return nil
	})

	if err != nil {
		return nil, nil, domerrors.ParseIntError(err)
	}

	return &loginUserRes, &authTokensDto, nil
}
