package authservice

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/auth"
	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
	"github.com/Samu-Amy/Shokora/internal/store/user"
	"github.com/golang-jwt/jwt/v5"
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
	verificationTokens, err := service.createVerificationTokensWithRetries(ctx, user.Id, auth.EmailVerification) // TODO: controlla che le scadenze nel db siano giuste
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
		service.logger.Warnw("Error sending welcome email", "error", err)

		// Set email "error" in response
		registerUserRes.IsEmailSent = false
	}

	// ----- AUTH -----

	// Create Refresh Token (soft error)
	authTokensDto, err := service.createNewSessionAndRefreshToken(ctx, user.Id)
	if err != nil {
		registerUserRes.HasAuthError = true
		return registerUserRes, nil, nil
	}

	// Create Access Token (soft error)
	err = service.addJWTAccessToken(authTokensDto, user.Id)
	if err != nil {
		registerUserRes.HasAuthError = true
		return registerUserRes, nil, nil
	}

	service.logger.Info("User and Tokens created, Email sent successfully", "userId", user.Id)

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
	// TODO: usa ParseIntError in ogni err return

	// ----- USER -----

	// Hash password
	// hashedPassword, err := service.hashPassword(payload.Password)
	// if err != nil {
	// 	service.logger.Warnw("Error hashing password", "error", err)
	// 	return nil, nil, domerrors.ParseIntError(err)
	// }

	// Get user
	// user, err := service.getUser(ctx, email, hashedPassword) // TODO: fai metodo per ottenere user da email e verificare (password, is_verified, is_active (?))
	// if err != nil {
	// 	return nil, nil, domerrors.ParseIntError(err)
	// }

	// Create Response payload with UserRes built from user model
	// loginUserRes := payloads.NewLoginUserRes(payloads.ToUserRes(*user))

	// ----- VERIFICATION -----

	// TODO: se non verificato o 2fa -> crea verificationTokens ed invia email, altrimenti crea authTokens (modificare return per distinguere i due casi)

	// ----- AUTH -----

	// Create Refresh Token
	// authTokenDto, err := service.createNewRefreshToken(ctx, user.Id)
	// if err != nil {
	// 	return nil, nil, domerrors.ParseIntError(err) // TODO: controlla parsing errore
	// }

	// service.logger.Info("User and Tokens created, Email sent successfully", "userId", user.Id)

	// return loginUserRes, authTokenDto, nil

	// TODO: nel login fai anche delete di tutti i refresh token scaduti per quell'utente (o in generale?) - ottenere un l'ultimo token creato per ogni session_id (join con order by) e se è scaduto -> sessione scaduta (?)

	return nil, nil, nil // TODO: modifica (quello sopra)
}

/*
Log out the user
*/
func (service *AuthService) LogoutUser(ctx context.Context, userId int64) error {

	// Delete session

	// TODO: elimina cookies in hanlder

	return nil
}

/*
Takes Access and Refresh Tokens, verifies and updates them (if necessary)

Return:
  - *payloads.AuthTokensCheckDto: auth tokens data required so set auth cookies (AuthTokensDto) + IsAccessTokenValid and UserId
  - error: domerrors (safe to send to the frontend)
*/
func (service *AuthService) HandleAuthTokensCheck(ctx context.Context, accessToken, plainRefreshToken string) (*payloads.AuthTokensCheckDto, error) {

	authTokensCheckDto := payloads.AuthTokensCheckDto{
		IsAccessTokenValid: false,
	}

	// Verify Access Token
	jwtToken, err := service.jwtAuthenticator.ValidateJWTToken(accessToken)
	if err == nil && jwtToken != nil {

		// Get user Id from claims
		claims := jwtToken.Claims.(jwt.MapClaims)

		userId, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err == nil {
			// Access Token valid -> set data and return
			authTokensCheckDto.IsAccessTokenValid = true
			authTokensCheckDto.UserId = userId

			return &authTokensCheckDto, nil
		}
	}

	// Access Token not valid -> Rotate Refresh Token
	hashedRefreshToken := auth.HashBase64Token(plainRefreshToken)
	authTokensDto, userId, err := service.rotateRefreshToken(ctx, hashedRefreshToken)
	if err != nil {
		return nil, domerrors.ParseIntError(err)
	}

	if userId == -1 {
		return nil, domerrors.ErrNotFound
	}

	// Create Access Token
	err = service.addJWTAccessToken(authTokensDto, userId)
	if err != nil {
		return nil, domerrors.ParseIntError(err)
	}

	// TODO: controlla che sia tutto giusto

	authTokensCheckDto.TokensDto = *authTokensDto

	// Create new access token

	return &authTokensCheckDto, nil
}
