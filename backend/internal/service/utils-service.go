package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/errorcodes"
	"github.com/Samu-Amy/Shokora/internal/store"
)

// - Timeouts -

const (
	regenerate_token_timeout = 10 * time.Second
)

// - Functions -

// Transaction wrapper
func withTransaction(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	// Create transaction
	transaction, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Defer rollback (in caso di panic)
	defer func() {
		if err != nil {
			_ = transaction.Rollback() // TODO: rollback error handling?
		}
	}()

	if err = fn(transaction); err != nil {
		return err
	}

	return transaction.Commit()
}

// - Verification Tokens -

func (service *AuthService) verifyOtp(ctx context.Context, verificationId int64, hashedOTP []byte, maxAttempts uint8, verificationType auth.VerificationType) (*store.OTPPayload, error) {

	// Get data
	otpQueryData, err := service.vTokensRepo.GetOTPData(ctx, verificationId, verificationType)
	if err != nil {
		switch {
		case errors.Is(err, errorcodes.ErrNotFound): // Not valid (id does not exists or wrong verificationType)
			return nil, errorcodes.ErrInvalid
		default:
			return nil, err // db/query error
		}
	}

	// Verify other data
	var verificationErr error

	if otpQueryData != nil && otpQueryData.Exp.Before(time.Now()) { // Check expiry
		verificationErr = errorcodes.InternalErrExpired

	} else if otpQueryData != nil && otpQueryData.Attempts >= maxAttempts { // Check attempts
		verificationErr = errorcodes.ErrMaxAttemptsExceeded
	}

	// Validate OTP // TODO: controlla
	validOtp := service.tokenAuthenticator.VerifyOTP(hashedOTP, otpQueryData.HashedOtp)
	if !validOtp {
		verificationErr = errorcodes.ErrInvalid

		// Increment attempts and Handle errors
		err = service.vTokensRepo.UpdateOtpAttempts(ctx, verificationId, maxAttempts)
		if err != nil && verificationErr != errorcodes.ErrInvalid {
			verificationErr = err // MaxAttemptsExceeded or db/query error
		}
		return nil, errorcodes.ErrInvalid
	}

	return otpQueryData, verificationErr
}
