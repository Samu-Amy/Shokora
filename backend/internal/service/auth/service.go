package authservice

import (
	"context"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/auth"
	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
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

	// TODO: aggiungi service.logger.Warnw("Error ...", "error", err) nei metodi del service "low level" usati qui

	// TODO: usa il parsing in ogni return

	// ----- USER -----

	// Hash password
	hashedPassword, err := service.hashPassword(payload.Password)
	if err != nil {
		service.logger.Warnw("Error hashing password", "error", err)
		return nil, nil, domerrors.ParseIntError(err)
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
		return nil, nil, domerrors.ParseIntError(err)
	}

	// Create Response payload with UserRes built from user model
	registerUserRes := payloads.NewRegisterUserRes(payloads.ToUserRes(*user))

	// ----- VERIFICATION -----

	// Create Email Verification Tokens (soft error)
	verificationTokens, err := service.createVerificationTokensWithRetries(ctx, user)
	if err != nil {
		registerUserRes.VerificationError = softerrors.SoftErrVerification
	}

	// Add verification id to payload
	registerUserRes.VerificationId = verificationId // TODO: ottieni id da createVerificationTokensWithRetries()

	// Send email (soft error)
	err = service.sendVerificationEmail(
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

	// ----- Auth -----

	// Create Refresh Token (soft error)
	authTokenDto, err := service.createNewRefreshToken(ctx, user.Id)
	if err != nil {
		resPayload.AuthError = true
	}

	// app.logger.Info("User and Tokens created, Email sent successfully")

	return registerUserRes, authTokenDto, nil
}
