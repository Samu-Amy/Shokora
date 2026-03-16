package authservice

import (
	"context"
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/auth"
	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
)

// ----- VERIFY EMAIL  -----

func (service *AuthService) VerifyEmailWithMagicLink(ctx context.Context, plainToken string) error {
	// TODO: usare FOR UPDATE nel get (per l'eliminazione)? usare transaction?

	err := service.txManager.WithTx(ctx, func(tx *sql.Tx) error {

		// Hash token
		hashedToken := auth.HashBase64Token(plainToken)

		// Verify and get data
		magicLinkTokenQueryData, err := service.vTokenRepo.GetValidMagicLinkData(ctx, tx, hashedToken, auth.EmailVerification)
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
		if err = service.vTokenRepo.Delete(ctx, tx, magicLinkTokenQueryData.VerificationId); err != nil { // If it fails to delete there are no problems
			service.logger.Errorw("failed deleting verification token", "error", err)
		}

		return nil
	})

	return err
}

func (service *AuthService) VerifyEmailWithOTP(ctx context.Context, payload payloads.OTPVerificationReq) error {

	err := service.txManager.WithTx(ctx, func(tx *sql.Tx) error {

		// Hash OTP
		hashedOTP := service.tokenAuthenticator.HashOTP(payload.OTP, auth.EmailVerification)

		// Verify and get data
		userId, err := service.verifyOtp(ctx, tx, payload.VerificationId, hashedOTP, service.config.Auth.OTP.MaxAttempts, auth.EmailVerification)
		if err != nil {
			service.logger.Warnw("Error verifying Otp", "error", err)
			return domerrors.ParseIntError(err)
		}

		// Verify user
		err = service.userRepo.SetIsVerified(ctx, userId)
		if err != nil {
			service.logger.Warnw("Error setting User is_verified", "error", err)
			return domerrors.ParseIntError(err)
		}

		// Delete token
		if err = service.vTokenRepo.Delete(ctx, tx, payload.VerificationId); err != nil { // If it fails to delete there are no problems
			service.logger.Errorw("failed deleting verification token", "error", err)
		}

		return nil
	})

	return err
}

// ----- PASSWORD RESET  -----

func (service *AuthService) ResetPasswordWithMagicLink(ctx context.Context, plainToken string) (string, error) {

	var resetSessionToken string

	err := service.txManager.WithTx(ctx, func(tx *sql.Tx) error {

		// Hash token
		hashedToken := auth.HashBase64Token(plainToken)

		// Verify and get data
		magicLinkTokenQueryData, err := service.vTokenRepo.GetValidMagicLinkData(ctx, tx, hashedToken, auth.PasswordReset)
		if err != nil {
			service.logger.Warnw("Error getting magic link Token", "error", err)
			return domerrors.ParseIntError(err)
		}

		// TODO: crea reset session token (token univoco di 32 Bytes come il magic link) - imposta scadenza (10min) in app e tokenAuthenticator (?)

		// Delete token
		if err = service.vTokenRepo.Delete(ctx, tx, magicLinkTokenQueryData.VerificationId); err != nil { // If it fails to delete there are no problems
			service.logger.Errorw("failed deleting verification token", "error", err)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return resetSessionToken, nil
}

func (service *AuthService) ResetPasswordWithOTP(ctx context.Context, payload payloads.OTPVerificationReq) (string, error) {

	var resetSessionToken string

	err := service.txManager.WithTx(ctx, func(tx *sql.Tx) error {

		// Hash OTP
		hashedOTP := service.tokenAuthenticator.HashOTP(payload.OTP, auth.PasswordReset)

		// Verify and get data
		userId, err := service.verifyOtp(ctx, tx, payload.VerificationId, hashedOTP, service.config.Auth.OTP.MaxAttempts, auth.PasswordReset)
		if err != nil {
			service.logger.Warnw("Error verifying Otp", "error", err)
			return domerrors.ParseIntError(err)
		}

		// TODO: crea reset session token (token univoco di 32 Bytes come il magic link) - imposta scadenza (10min) in app e tokenAuthenticator (?)

		// Delete token
		if err = service.vTokenRepo.Delete(ctx, tx, payload.VerificationId); err != nil { // If it fails to delete there are no problems
			service.logger.Errorw("failed deleting verification token", "error", err)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return resetSessionToken, nil
}

// ----- TWO FACTOR AUTH  -----

func (service *AuthService) TwoFactorAuthWithOTP(ctx context.Context, payload payloads.OTPVerificationReq) (*payloads.AuthTokensDto, error) {

	var authTokensDto *payloads.AuthTokensDto
	var err error

	err = service.txManager.WithTx(ctx, func(tx *sql.Tx) error {

		// Hash OTP
		hashedOTP := service.tokenAuthenticator.HashOTP(payload.OTP, auth.TwoFactorAuth)

		// Verify and get data
		userId, err := service.verifyOtp(ctx, tx, payload.VerificationId, hashedOTP, service.config.Auth.OTP.MaxAttempts, auth.TwoFactorAuth)
		if err != nil {
			service.logger.Warnw("Error verifying Otp", "error", err)
			return domerrors.ParseIntError(err)
		}

		// Delete old sessions
		_ = service.userSessionRepo.DeleteExpired(ctx, userId)

		// Create Auth Tokens
		authTokensDto, err = service.createNewAuthTokens(ctx, userId)
		if err != nil {
			return domerrors.ParseIntError(err)
		}

		// Delete token
		if err = service.vTokenRepo.Delete(ctx, tx, payload.VerificationId); err != nil { // If it fails to delete there are no problems
			service.logger.Errorw("failed deleting verification token", "error", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return authTokensDto, nil
}
