package authservice

import (
	"context"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/auth"
	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
)

// ----- VERIFY EMAIL  -----

/*
Errors
  - ErrInvalid
  - Other db errors
*/
func (service *AuthService) VerifyEmailWithToken(ctx context.Context, plainToken string) error {
	// TODO: usare FOR UPDATE nel get (per l'eliminazione)? usare transaction?

	// Hash token
	hashedToken := auth.HashBase64Token(plainToken)

	// Verify and Get data
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

/*
Errors
  - ErrInvalid
  - InternalErrExpired
  - ErrMaxAttemptsExceeded
  - Other db errors
*/
func (service *AuthService) VerifyEmailWithOTP(ctx context.Context, payload payloads.OTPVerificationReq) error {

	// Hash OTP
	hashedOTP := service.tokenAuthenticator.HashOTP(payload.OTP, auth.EmailVerification)

	// Get data
	otpQueryData, err := service.verifyOtp(ctx, payload.VerificationId, hashedOTP, service.config.Auth.OTP.MaxAttempts, auth.EmailVerification)
	if err != nil {
		service.logger.Warnw("Error getting Otp", "error", err)
		return err
	}

	// Verify user
	err = service.userRepo.SetIsVerified(ctx, otpQueryData.UserId)
	if err != nil {
		service.logger.Warnw("Error setting User is_verified", "error", err)
		return domerrors.ParseIntError(err)
	}

	// Delete token
	_ = service.vTokenRepo.Delete(ctx, payload.VerificationId) // If it fails to delete there are no problems

	return nil
}
