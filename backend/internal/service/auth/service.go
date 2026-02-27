package authservice

import (
	"context"
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/auth"
	softerrors "github.com/Samu-Amy/Shokora/internal/errors/soft"
	"github.com/Samu-Amy/Shokora/internal/store/user"
)

/*
Creates a new user account and manages verification and authentication:
  - create user
  - create verification tokens (magic link and otp)
  - send email with verification tokens
  - create auth (access and refresh) tokens

Return:
  - *payloads.RegisterUserRes: response data (user, verification id and soft errors for verification and auth) to send to the frontend
  - *payloads.AuthTokensDto: auth token data required so set auth cookies
  - error: domerrors (safe to send to the frontend)
*/
func (service *AuthService) RegisterUser(ctx context.Context, payload payloads.RegisterUserReq) (*payloads.RegisterUserRes, *payloads.AuthTokensDto, error) {

	// TODO: ogni controllo sulle date da fare nel db (?)

	// Hash password
	hashedPassword, err := service.hashPassword(payload.Password)
	if err != nil {
		service.logger.Warnw("Error hashing password", "error", err)
		return nil, domerrors.ParseIntError(err) // TODO: usa il parsing in ogni return
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
		return nil, domerrors.ParseIntError(err)
	}

	// Build user payload from updated struct
	userRes := payloads.ToUserRes(*user)
	
	// Create Response Payload with user
	registerUserRes := payloads.NewRegisterUserRes(userRes)
	
	// -------------------------

	//

	// Generate Refresh Token

	// Generate Access Token

	// -------------------------

	// Create Refresh Token
	refreshToken, err := service.createNewRefreshToken(ctx, user.Id)
	if err != nil {
		service.logger.Warnw("Error generating refresh token", "error", err)
		resPayload.AuthError = true
	}

	// Generate verificationTokens (Magic Link and OTP)
	verificationTokens, err := app.tokenAuthenticator.CreateVerificationTokens(auth.EmailVerification)
	if err != nil {
		app.logger.Warnw("error generating verification tokens", "error", err)

		registerUserDto.RegisterUserRes.VerificationError = softerrors.SoftErrVerification // Add error to payload // TODO: verifica che la serializzazione funzioni correttamente

		//* Return user, verificationID and error
		// if err := app.jsonResponse(w, http.StatusCreated, registerUserRes); err != nil {
		// 	app.internalServerError(w, r, err)
		// }
		// return
	}

	// Create Email Verification Tokens
	verificationId, err := app.service.Auth.CreateVerificationTokensWithRetries(ctx, user, verificationTokens)
	if err != nil {
		app.logger.Warnw("error creating email verification tokens in db", "error", err)

		registerUserDto.RegisterUserRes.VerificationError = softerrors.SoftErrVerification // Add error to payload

		//* Return user, verificationID and error
		// if err := app.jsonResponse(w, http.StatusCreated, registerUserRes); err != nil {
		// 	app.internalServerError(w, r, err)
		// }
		// return
	}

	registerUserRes.VerificationId = verificationId // Ad verification id to payload

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

		registerUserDto.RegisterUserRes.VerificationError = softerrors.SoftErrEmailNotSent // Add error to payload

		//* Return user, verificationID and error
		// if err := app.jsonResponse(w, http.StatusCreated, registerUserRes); err != nil {
		// 	app.internalServerError(w, r, err)
		// }
		// return

		// TODO: dire di riprovare più tardi? -> l'utente può accedere ma non può ordinare (ha come opzioni di re-inviare la mail di verifica oppure eliminare l'account (e il token))
	}

	app.logger.Info("User and Tokens created, Email sent successfully")

	return registerUserRes, , nil
}
