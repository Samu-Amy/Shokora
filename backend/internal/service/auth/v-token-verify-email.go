package authservice

import (
	"context"
	"errors"

	"github.com/Samu-Amy/Shokora/internal/auth"
	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
)

// ----- VERIFY EMAIL  -----

/*
Errors
  - ErrInvalid
  - Other db errors
*/
func (service *AuthService) verifyEmailWithToken(ctx context.Context, hashedToken []byte) error {

	// Verify and Get data
	magicLinkTokenQueryData, err := service.vTokenRepo.GetValidMagicLinkData(ctx, hashedToken, auth.EmailVerification)
	if err != nil {
		// log.Printf("Verify OTP Error: %v", err)
		switch {
		case errors.Is(err, domerrors.ErrNotFound): // Token not valid // TODO: si dovrebbero ottenere solo interrors
			return domerrors.ErrInvalid
		default:
			return err
		}
	}

	// Verify user
	err = service.userRepo.Verify(ctx, magicLinkTokenQueryData.UserId)
	if err != nil {
		// log.Printf("Verify User Error: %v", err)
		return err
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
func (service *AuthService) verifyEmailWithOTP(ctx context.Context, verificationId int64, hashedOTP []byte, maxAttempts uint8) error {

	// Get data
	otpQueryData, err := service.verifyOtp(ctx, verificationId, hashedOTP, maxAttempts, auth.EmailVerification)
	if err != nil {
		// log.Printf("Verify OTP Error: %v", err)
		return err
	}

	// Verify user
	err = service.userRepo.Verify(ctx, otpQueryData.UserId)
	if err != nil {
		// log.Printf("Verify User Error: %v", err)
		return err
	}

	// Delete token
	_ = service.vTokenRepo.Delete(ctx, verificationId) // If it fails to delete there are no problems

	return nil
}
