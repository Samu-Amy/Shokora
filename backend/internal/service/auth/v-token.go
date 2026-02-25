package authservice

import (
	"context"
	"errors"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/errorcodes"
	"github.com/Samu-Amy/Shokora/internal/store"
	"github.com/Samu-Amy/Shokora/internal/store/user"
)

// ----- CREATE VERIFICATION TOKENS -----

func (service *AuthService) CreateVerificationTokensWithRetries(ctx context.Context, user *user.User, verificationTokens *auth.VerificationTokens) (*int64, error) {

	ctxWithTimeout, cancel := context.WithTimeout(ctx, auth.REGENERATE_TOKEN_TIMEOUT)
	defer cancel()

	// Create Tokens in db
	verificationId, err := service.vTokenRepo.CreateTokens(ctxWithTimeout, user.Id, verificationTokens)
	if err == nil {
		return verificationId, nil // OK, return no error

	} else if !errors.Is(err, errorcodes.InternalErrDuplicateToken) {
		return nil, err // Error (can't retry)
	}

	// Retries
	for range service.tokenAuthenticator.MaxRetries - 1 {

		// Timeout
		select {
		case <-ctxWithTimeout.Done():
			return nil, ctxWithTimeout.Err()
		default:
		}

		// Regenerate Tokens (if error is "duplicate token")
		switch {

		// Duplicated Magic Link Token
		case errors.Is(err, errorcodes.InternalErrDuplicateToken):
			err = service.tokenAuthenticator.RegenerateMagicLinkToken(verificationTokens)
			if err != nil {
				continue // skip iteration
			}

			verificationId, err = service.vTokenRepo.CreateTokens(ctx, user.Id, verificationTokens)
			if err == nil {
				return verificationId, nil // OK, return no error
			}

		default:
			return nil, err // Error is not solvable (not "duplicate token") -> return it
		}
	}

	return nil, errorcodes.ErrMaxRetriesExceeded // Couldn't regenerate and save token successfully -> return error "max retries exceeded"
}

// ----- VERIFY EMAIL  -----

/*
Errors
  - ErrInvalid
  - Other db errors
*/
func (service *AuthService) VerifyEmailWithToken(ctx context.Context, hashedToken []byte) error {

	// Verify and Get data
	magicLinkTokenQueryData, err := service.vTokenRepo.VerifyMagicLink(ctx, hashedToken, auth.EmailVerification)
	if err != nil {
		// log.Printf("Verify OTP Error: %v", err)
		switch {
		case errors.Is(err, errorcodes.ErrNotFound): // Token not valid
			return errorcodes.ErrInvalid
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
func (service *AuthService) VerifyEmailWithOTP(ctx context.Context, verificationId int64, hashedOTP []byte, maxAttempts uint8) error {

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

// ----- RESET PASSWORD  -----

// ----- TWO FACTOR AUTH  -----

// - Verification Tokens -

func (service *AuthService) verifyOtp(ctx context.Context, verificationId int64, hashedOTP []byte, maxAttempts uint8, verificationType auth.VerificationType) (*store.OTPPayload, error) {
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
