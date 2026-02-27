package authservice

import (
	"context"
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/auth"
	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
	"github.com/Samu-Amy/Shokora/internal/store/user"
)

/*
Executes the user registration:
  - create new user in db
    -
*/
func (service *AuthService) RegisterUser(ctx context.Context, payload payloads.RegisterUserReq) (*payloads.RegisterUserRes, error) {
	// TODO: ritorna anche i dati per i cookies

	// Hash password
	hashedPassword, err := service.hashPassword(payload.Password)
	if err != nil {
		service.logger.Warnw("Error hashing password", "error", err)
		return nil, err
	}

	// Build user struct from payload data
	user := &user.User{
		FirstName:    payload.FirstName,
		LastName:     payload.LastName,
		Email:        payload.Email,
		PasswordHash: hashedPassword,
		ImageUrl:     payload.ImageUrl,
		BirthDate:    payload.BirthDate,
	}

	// Create user in db and update its struct
	if err := service.createUser(ctx, user); err != nil {
		service.logger.Warnw("Error creating user", "error", err)
		return nil, err
	}

	// Build user payload from updated struct
	userRes := payloads.ToUserRes(*user)

	// -------------------------

	//

	//

	// -------------------------

	// Create Response Payload with user
	registerUserRes := &payloads.RegisterUserRes{
		User: userRes,
	}

	// TODO: continua

	// Create Refresh Token
	refreshToken, err := service.generateNewRefreshToken(ctx, user.Id)
	if err != nil {
		service.logger.Warnw("Error generating refresh token", "error", err)
		resPayload.AuthError = true
	}

	// Auth cookies
	err = app.setAuthCookies(w, user.Id, refreshToken.PlainToken, refreshToken.ExpiresAt)
	if err != nil {
		resPayload.AuthError = true
	}

	// Generate verificationTokens (Magic Link and OTP)
	verificationTokens, err := app.tokenAuthenticator.CreateVerificationTokens(auth.EmailVerification)
	if err != nil {
		app.logger.Warnw("error generating verification tokens", "error", err)

		resPayload.VerificationError = domerrors.ErrVerification.Error() // Add error to payload

		//* Return user, verificationID and error
		if err := app.jsonResponse(w, http.StatusCreated, resPayload); err != nil {
			app.internalServerError(w, r, err)
		}
		return
	}

	// Create Email Verification Tokens
	verificationId, err := app.service.Auth.CreateVerificationTokensWithRetries(ctx, user, verificationTokens)
	if err != nil {
		app.logger.Warnw("error creating email verification tokens in db", "error", err)

		resPayload.VerificationError = domerrors.ErrVerification.Error() // Add error to payload

		//* Return user, verificationID and error
		if err := app.jsonResponse(w, http.StatusCreated, resPayload); err != nil {
			app.internalServerError(w, r, err)
		}
		return
	}

	resPayload.VerificationId = verificationId // Ad verification id to payload

	// Send email
	err = app.SendVerificationEmail(
		ctx,
		auth.EmailVerification,
		user.FirstName,
		user.Email,
		verificationTokens.PlainMagicLinkToken,
		verificationTokens.PlainOTP,
		verificationTokens.MagicLinkTokenExp,
		verificationTokens.OTPExp,
	)
	if err != nil {
		app.logger.Warnw("error sending welcome email", "error", err)

		resPayload.VerificationError = domerrors.ErrEmailNotSent.Error() // Add error to payload

		//* Return user, verificationID and error
		if err := app.jsonResponse(w, http.StatusCreated, resPayload); err != nil {
			app.internalServerError(w, r, err)
		}
		return

		// TODO: dire di riprovare più tardi? -> l'utente può accedere ma non può ordinare (ha come opzioni di re-inviare la mail di verifica oppure eliminare l'account (e il token))
	}

	app.logger.Info("User and Tokens created, Email sent successfully")

	return resPayload, nil
}
