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

	var user *user_repo.User
	var loginUserRes payloads.LoginUserRes
	var authTokensDto *payloads.AuthTokensDto
	var err error

	err = service.txManager.WithTx(ctx, func(tx *sql.Tx) error {

		// ----- OAUTH TOKEN -----

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
		firstName, firstNameOk := tokenPayload.Claims["given_name"].(string)
		lastName, lastNameOk := tokenPayload.Claims["family_name"].(string)
		emailVerified, emailVerifiedOk := tokenPayload.Claims["email_verified"].(bool)
		// picture, _ := tokenPayload.Claims["picture"].(string)

		// ----- USER -----

		// Get user by googleId (if exists, login)
		user, err = service.userRepo.GetByGoogleId(ctx, googleId)
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

				// User does not exists -> Create user
				if !firstNameOk {
					return domerrors.ErrInvalid
				}

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

				// TODO: non c'è la password -> come fare (va bene lasciarla null se c'è google id, ma bloccando login normale? -> la si può in qualche modo cambiare dopo (però l'update richiede la vecchia password, che però non esiste))?

				err = service.userRepo.Create(ctx, tx, user)
				if err != nil {
					return err
				}
			}

			// Check if verified (if not -> error)
			if !user.IsVerified {
				return domerrors.ErrNotVerified // TODO: ricorda di dire nel FRONTEND che serve avere account verificato se si vuole accedere ad un account email e password usando google
			}

			// Add google_id (connect email/password and google accounts)
			err = service.userRepo.SetGoogleId(ctx, tx, user.Id, googleId) // soft error
			if err != nil {
				service.logger.Warnw("Error setting googleId to user ", "error", err, "userId", user.Id)
			}
		}

		return nil
	})

	if err != nil {
		return nil, nil, domerrors.ParseIntError(err)
	}

	// ----- VERIFICATION -----

	// Check if 2FA verification is needed
	hasTwoFactorAuth, err := service.userSettingsRepo.GetHasTwoFactorAuthById(ctx, user.Id)
	if err != nil {
		service.logger.Warnw("Error getting hasTwoFactorAuth", "error", err)
		return nil, nil, domerrors.ParseIntError(err)
	}

	if hasTwoFactorAuth {
		//  - 2FA needed -

		// Create Verification Tokens
		verificationTokens, err := service.createVerificationTokensWithRetries(ctx, user.Id, auth.TwoFactorAuth)
		if err != nil {
			service.logger.Warnw("Error creating 2FA verification token ", "error", err)
			return nil, nil, domerrors.ParseIntError(err)
		}

		// Send email (soft error)
		if verificationTokens != nil {

			// Save verification id in response
			loginUserRes.VerificationId = &verificationTokens.VerificationId //* If registerUserRes.VerificationId == nil -> error during verification (tokens not created)

			err = service.sendVerificationEmail(
				ctx,
				auth.TwoFactorAuth,
				user.FirstName,
				user.Email,
				verificationTokens.PlainMagicLinkToken,
				verificationTokens.PlainOTP,
			)
			if err != nil {
				service.logger.Warnw("Error sending verification email", "error", err)

				// Set email "error" in response
				loginUserRes.IsEmailSent = false
			} else {
				loginUserRes.IsEmailSent = true
			}

			// service.logger.Info("Verification email sent", "userId", user.Id, "verificationType", auth.TwoFactorAuth)
		}
	} else {
		// Add user to payload (no 2FA required)
		loginUserRes.User = payloads.ToUserRes(*user)

		// ----- AUTH TOKENS -----

		// Delete old sessions
		_ = service.userSessionRepo.DeleteExpired(ctx, user.Id)

		// Create Auth Tokens
		authTokensDto, err = service.createNewAuthTokens(ctx, user.Id)
		if err != nil {
			return nil, nil, domerrors.ParseIntError(err)
		}

		// service.logger.Info("User logged, Tokens created", "userId", user.Id)
	}

	return &loginUserRes, authTokensDto, nil
}
