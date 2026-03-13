package authservice

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
	interrors "github.com/Samu-Amy/Shokora/internal/errors/int"
	vtoken "github.com/Samu-Amy/Shokora/internal/store/verification-token"
	"github.com/google/uuid"
)

// ----- VERIFY OTP -----

func (service *AuthService) verifyOtp(ctx context.Context, verificationId uuid.UUID, hashedOTP []byte, maxAttempts uint8, verificationType auth.VerificationType) (*vtoken.OTPVerificationData, error) {

	var otpQueryData *vtoken.OTPVerificationData
	var err error

	err = service.txManager.WithTx(ctx, func(tx *sql.Tx) error {

		// Get data
		otpQueryData, err = service.vTokenRepo.GetOtpData(ctx, tx, verificationId, verificationType)
		if err != nil {
			service.logger.Warnw("Error getting otp data", "error", err, "verificationId", verificationId)
			// Not valid (id does not exists or wrong verificationType)
			if errors.Is(err, interrors.IErrNotFound) {
				return interrors.IErrInvalid
			}

			return err // db/query error
		}

		// Verify attempts
		if otpQueryData.Attempts >= maxAttempts {
			// service.logger.Warn("Max attempts for otp")
			return interrors.IErrMaxRetriesExceeded
		}

		// Verify expiry
		if otpQueryData.ExpiresAt.Before(time.Now()) {
			// service.logger.Warn("Expired otp")
			return interrors.IErrExpired
		}

		// Validate OTP
		isOtpValid := service.tokenAuthenticator.VerifyOTP(hashedOTP, otpQueryData.HashedOtp)
		if !isOtpValid {

			// Increment attempts and handle errors
			err = service.vTokenRepo.IncrementOtpAttempts(ctx, tx, verificationId, maxAttempts)
			if err != nil {
				service.logger.Warnw("Error updating otp attempts", "error", err, "verificationId", verificationId)
			}

			// Attempts updated successfully but OTP not valid
			if errors.Is(err, interrors.IErrNoRowsAffected) { // Max attempts exceeded
				err = interrors.IErrMaxAttemptsExceeded
			}

			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return otpQueryData, nil
}
