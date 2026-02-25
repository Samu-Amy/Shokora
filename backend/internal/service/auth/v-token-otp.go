package authservice

import (
	"context"
	"errors"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/errorcodes"
	vtoken "github.com/Samu-Amy/Shokora/internal/store/verification-token"
)

// ----- VERIFY OTP -----

func (service *AuthService) verifyOtp(ctx context.Context, verificationId int64, hashedOTP []byte, maxAttempts uint8, verificationType auth.VerificationType) (*vtoken.OTPPayload, error) {
	// TODO: usare transaction (ed usare FOR UPDATE nel get?) per GetOtpData e UpdateOtpAttempts?

	// Get data
	otpQueryData, err := service.vTokenRepo.GetOtpData(ctx, verificationId, verificationType)
	if err != nil {
		switch {
		case errors.Is(err, errorcodes.ErrNotFound): // Not valid (id does not exists or wrong verificationType)
			return nil, errorcodes.ErrInvalid
		default:
			return nil, err // db/query error
		}
	}

	// Verify attempts
	if otpQueryData.Attempts >= maxAttempts {
		return nil, errorcodes.ErrMaxAttemptsExceeded
	}

	// Verify expiry
	if otpQueryData.Exp.Before(time.Now()) {
		return nil, errorcodes.InternalErrExpired
	}

	// Validate OTP
	isOtpValid := service.tokenAuthenticator.VerifyOTP(hashedOTP, otpQueryData.HashedOtp)
	if !isOtpValid {

		// Increment attempts and Handle errors
		err = service.vTokenRepo.UpdateOtpAttempts(ctx, verificationId, maxAttempts)

		// Attempts updated successfully but OTP not valid
		if err == nil {
			switch {
			case errors.Is(err, errorcodes.ErrNotFound):
				err = errorcodes.ErrInvalid // VerificationId is not valid
			case errors.Is(err, errorcodes.InternalErrNoRowsAffected): // Max attempts exceeded
				err = errorcodes.ErrMaxAttemptsExceeded
			default:
				err = errorcodes.ErrInvalid
			}
		}

		return nil, err
	}

	return otpQueryData, nil
}
