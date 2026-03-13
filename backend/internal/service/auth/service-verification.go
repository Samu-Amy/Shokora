package authservice

import (
	"context"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/auth"
	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
)

// ----- VERIFY EMAIL  -----

func (service *AuthService) VerifyEmailWithMagicLink(ctx context.Context, plainToken string) error {
	// TODO: usare FOR UPDATE nel get (per l'eliminazione)? usare transaction?

	// Hash token
	hashedToken := auth.HashBase64Token(plainToken)

	// Verify and get data
	magicLinkTokenQueryData, err := service.vTokenRepo.GetValidMagicLinkData(ctx, hashedToken, auth.EmailVerification)
	if err != nil {
		service.logger.Warnw("Error getting magic link Token", "error", err)
		return domerrors.ParseIntError(err)
	}

	// Verify user
	err = service.userRepo.SetIsVerified(ctx, magicLinkTokenQueryData.UserId)
	if err != nil {
		service.logger.Warnw("Error setting User is_verified", "error", err)
		return domerrors.ParseIntError(err)
	}

	// Delete token
	_ = service.vTokenRepo.Delete(ctx, magicLinkTokenQueryData.VerificationId) // If it fails to delete there are no problems

	return nil
}

func (service *AuthService) VerifyEmailWithOTP(ctx context.Context, payload payloads.OTPVerificationReq) error {

	// Hash OTP
	hashedOTP := service.tokenAuthenticator.HashOTP(payload.OTP, auth.EmailVerification)

	// Verify and get data
	userId, err := service.verifyOtp(ctx, payload.VerificationId, hashedOTP, service.config.Auth.OTP.MaxAttempts, auth.EmailVerification)
	if err != nil {
		service.logger.Warnw("Error getting OTP", "error", err)
		return domerrors.ParseIntError(err)
	}

	// Verify user
	err = service.userRepo.SetIsVerified(ctx, userId)
	if err != nil {
		service.logger.Warnw("Error setting User is_verified", "error", err)
		return domerrors.ParseIntError(err)
	}

	// Delete token
	_ = service.vTokenRepo.Delete(ctx, payload.VerificationId) // If it fails to delete there are no problems

	return nil
}

// ----- PASSWORD RESET  -----

func (service *AuthService) ResetPasswordWithMagicLink(ctx context.Context, plainToken string) (string, error) {
	// TODO: per reset password: crea v-tokens -> verifica tokens e crea un reset session token (token univoco di 32 Bytes come il magic link)

	// Hash token
	hashedToken := auth.HashBase64Token(plainToken)

	// Verify and get data
	magicLinkTokenQueryData, err := service.vTokenRepo.GetValidMagicLinkData(ctx, hashedToken, auth.EmailVerification)
	if err != nil {
		service.logger.Warnw("Error getting magic link Token", "error", err)
		return "", domerrors.ParseIntError(err)
	}

	var resetSessionToken string

	return resetSessionToken, nil
}

func (service *AuthService) ResetPasswordWithOTP(ctx context.Context, payload payloads.OTPVerificationReq) (string, error) {

	// Hash OTP
	hashedOTP := service.tokenAuthenticator.HashOTP(payload.OTP, auth.EmailVerification)

	// Verify and get data
	userId, err := service.verifyOtp(ctx, payload.VerificationId, hashedOTP, service.config.Auth.OTP.MaxAttempts, auth.EmailVerification)
	if err != nil {
		service.logger.Warnw("Error getting Otp", "error", err)
		return "", domerrors.ParseIntError(err)
	}
	var resetSessionToken string

	return resetSessionToken, nil
}

// ----- TWO FACTOR AUTH  -----

func (service *AuthService) TwoFactorAuthWithOTP(ctx context.Context, payload payloads.OTPVerificationReq) (*payloads.AuthTokensDto, error) {

	// Hash OTP
	hashedOTP := service.tokenAuthenticator.HashOTP(payload.OTP, auth.EmailVerification)

	// Verify and get data
	userId, err := service.verifyOtp(ctx, payload.VerificationId, hashedOTP, service.config.Auth.OTP.MaxAttempts, auth.EmailVerification)
	if err != nil {
		service.logger.Warnw("Error getting Otp", "error", err)
		return nil, domerrors.ParseIntError(err)
	}

	// Delete old sessions
	_ = service.userSessionRepo.DeleteExpired(ctx, userId)

	// Create Auth Tokens
	authTokensDto, err := service.createAuthTokens(ctx, userId)
	if err != nil {
		return nil, domerrors.ParseIntError(err)
	}

	return authTokensDto, nil
}
