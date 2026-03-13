package authservice

import (
	"context"
	"errors"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/auth"
	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
	interrors "github.com/Samu-Amy/Shokora/internal/errors/int"
	"github.com/Samu-Amy/Shokora/internal/store/user"
)

/*
Creates a new user account and manages verification and authentication:
  - create user
  - create verification tokens (magic link and otp)
  - send email with verification tokens
  - create auth (access and refresh) tokens

Returns:
  - *payloads.RegisterUserRes: response data (user, verification id and soft errors for verification and auth) to send to the frontend
  - *payloads.AuthTokensDto: auth tokens data required so set auth cookies
  - error: domerrors (safe to send to the frontend)
*/
func (service *AuthService) RegisterUser(ctx context.Context, payload payloads.RegisterUserReq) (*payloads.RegisterUserRes, *payloads.AuthTokensDto, error) {

	// ----- USER -----

	// Hash password
	hashedPassword, err := service.hashPassword(payload.Password)
	if err != nil {
		service.logger.Warnw("Error hashing password", "error", err)
		return nil, nil, domerrors.ParseIntError(err)
	}

	// Build User struct from payload data
	user := &user.User{
		FirstName:    payload.FirstName,
		LastName:     payload.LastName,
		Email:        payload.Email,
		PasswordHash: hashedPassword,
		ImageUrl:     payload.ImageUrl,
		BirthDate:    payload.BirthDate,
	}

	// Create User in db and update its struct
	if err := service.createUser(ctx, user); err != nil {
		return nil, nil, domerrors.ParseIntError(err)
	}

	// Create Response payload with UserRes built from user model
	registerUserRes := payloads.NewRegisterUserRes(payloads.ToUserRes(*user))

	// ----- VERIFICATION -----

	// Create Email Verification Tokens (soft error)
	verificationTokens, err := service.createVerificationTokensWithRetries(ctx, user.Id, auth.EmailVerification)
	if err == nil {
		// Add verification id to response
		registerUserRes.VerificationId = &verificationTokens.VerificationId //* If registerUserRes.VerificationId == nil -> error during verification (tokens not created)
	}

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
		service.logger.Warnw("Error sending verification email", "error", err)

		// Set email "error" in response
		registerUserRes.IsEmailSent = false
	}

	// ----- AUTH -----

	// Create Auth Tokens (soft error)
	authTokensDto, err := service.createAuthTokens(ctx, user.Id)
	if err != nil {
		registerUserRes.HasAuthError = true

		service.logger.Info("Only User created", "userId", user.Id)

		return registerUserRes, nil, nil
	}

	service.logger.Info("User and Tokens created", "userId", user.Id)

	return registerUserRes, authTokensDto, nil
}

/*
Get the user account and manages verification and authentication:
  - get user
  - check (es. password, is_verified)
  - if not verified or has 2fa active create verification tokens (magic link and/or otp) and send them
  - if verified and doesn't have 2fa active create auth (access and refresh) tokens

Returns:
  - *payloads.LoginUserRes: response data to send to the frontend
  - *payloads.AuthTokensDto: auth tokens data required so set auth cookies
  - error: domerrors (safe to send to the frontend)
*/
func (service *AuthService) LoginUser(ctx context.Context, payload payloads.LoginUserReq) (*payloads.LoginUserRes, *payloads.AuthTokensDto, error) {

	loginUserRes := payloads.LoginUserRes{}
	var verificationType auth.VerificationType
	isVerificationRequired := false

	// ----- USER -----

	// Get user
	user, err := service.getUser(ctx, payload.Email, []byte(payload.Password)) // TODO: aggiungi controllo compleanno (per non farlo lato frontend) e two factor auth check (magari join con settings table?)
	if err != nil {

		if errors.Is(err, interrors.IErrNotVerified) {

			// Not verified
			verificationType = auth.EmailVerification
			isVerificationRequired = true

		} else if errors.Is(err, interrors.IErrTwoFactorAuthReqired) {

			// 2FA required
			verificationType = auth.TwoFactorAuth
			isVerificationRequired = true

		} else {
			return nil, nil, domerrors.ParseIntError(err)
		}
	}

	// If no 2fa required
	if !(isVerificationRequired && verificationType == auth.TwoFactorAuth) {

		// Create Response payload with UserRes built from user model
		userRes := payloads.ToUserRes(*user)
		loginUserRes.User = &userRes
	}

	// ----- VERIFICATION -----

	if isVerificationRequired {

		// Create Verification Tokens (soft error)
		verificationTokens, err := service.createVerificationTokensWithRetries(ctx, user.Id, verificationType)
		if err == nil {
			// Add verification id to response
			loginUserRes.VerificationId = &verificationTokens.VerificationId //* If registerUserRes.VerificationId == nil -> error during verification (tokens not created)
		}

		// Send email (soft error)
		if verificationTokens != nil {
			err = service.sendVerificationEmail(
				ctx,
				verificationType,
				user.FirstName,
				user.Email,
				verificationTokens.PlainMagicLinkToken,
				verificationTokens.PlainOTP,
				verificationTokens.MagicLinkTokenExp,
				verificationTokens.OTPExp,
			)
			if err != nil {
				service.logger.Warnw("Error sending verification email", "error", err)

				// Set email "error" in response
				loginUserRes.IsEmailSent = false
			}

			service.logger.Info("Verification email sent", "userId", user.Id, "verificationType", verificationType)
		}
	}

	// ----- AUTH -----

	var authTokensDto *payloads.AuthTokensDto

	// If no 2fa required
	if !(isVerificationRequired && verificationType == auth.TwoFactorAuth) {

		// Delete old sessions
		_ = service.userSessionRepo.DeleteExpired(ctx, user.Id)

		// Create Auth Tokens
		authTokensDto, err = service.createAuthTokens(ctx, user.Id)
		if err != nil {
			return nil, nil, domerrors.ParseIntError(err) // TODO: controlla parsing errore
		}

		service.logger.Info("User logged, Tokens created", "userId", user.Id)
	}

	return &loginUserRes, authTokensDto, nil
}

/*
Log out the user
*/
func (service *AuthService) LogoutUser(ctx context.Context, sessionId int64) error {

	// Delete session
	err := service.userSessionRepo.Delete(ctx, sessionId)
	if err != nil {
		service.logger.Warnw("Error deleting Session", "error", err)
		return domerrors.ParseIntError(err)
	}

	return nil
}

/*
Takes Access and Refresh Tokens, verifies and updates them (if necessary)

Return:
  - *payloads.AuthTokensCheckDto: auth tokens data required so set auth cookies (AuthTokensDto) + IsAccessTokenValid, UserId and SessionId
  - error: domerrors (safe to send to the frontend)
*/
func (service *AuthService) HandleAuthTokensCheck(ctx context.Context, accessToken, plainRefreshToken string) (*payloads.AuthTokensCheckDto, error) {

	// Verify Access Token
	authTokensCheckDto, err := service.checkAccessToken(accessToken)
	if err == nil {
		return authTokensCheckDto, nil // Access Token valid -> return early
	} else if !errors.Is(err, interrors.IErrExpired) {
		return nil, domerrors.ErrUnauthorized // Tokens doesn't correspond -> something is wrong
	}

	// Rotate Refresh Token (Access Token not valid)
	hashedRefreshToken := auth.HashBase64Token(plainRefreshToken)

	authTokensCheckDto, err = service.rotateRefreshToken(ctx, hashedRefreshToken)
	if err != nil {
		return nil, domerrors.ParseIntError(err)
	}

	// Create new Access Token (and update authTokensDto)
	err = service.addJWTAccessToken(&authTokensCheckDto.TokensDto, authTokensCheckDto.SessionId, authTokensCheckDto.UserId) // TODO: aggiungere tokenId in jwt (UserClaims)
	if err != nil {
		return nil, domerrors.ParseIntError(err)
	}

	return authTokensCheckDto, nil
}
