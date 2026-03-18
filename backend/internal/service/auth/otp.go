package authservice

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
	interrors "github.com/Samu-Amy/Shokora/internal/errors/int"
	"github.com/google/uuid"
)

// ----- VERIFY OTP -----

/*
Get the OTP from db, validate it and return userId
*/
func (service *AuthService) verifyOtp(ctx context.Context, tx *sql.Tx, verificationId uuid.UUID, hashedOTP []byte, maxAttempts uint8, verificationType auth.VerificationType) (int64, error) {

	// Get data
	otpQueryData, err := service.vTokenRepo.GetOtpData(ctx, tx, verificationId, verificationType)
	if err != nil {
		service.logger.Warnw("Error getting otp data", "error", err, "verificationId", verificationId)

		// Not valid (id does not exists or wrong verificationType)
		if errors.Is(err, interrors.IErrNotFound) {
			return 0, interrors.IErrInvalid
		}

		return 0, err // db/query error
	}

	// Verify attempts
	if otpQueryData.Attempts >= maxAttempts {
		service.logger.Warn("Max attempts for otp")
		return 0, interrors.IErrMaxRetriesExceeded
	}

	// Verify expiry
	if otpQueryData.ExpiresAt.Before(time.Now().UTC()) {
		service.logger.Warn("Expired otp")
		return 0, interrors.IErrExpired
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
			return 0, interrors.IErrMaxAttemptsExceeded
		}

		return 0, err
	}

	return otpQueryData.UserId, nil
}
