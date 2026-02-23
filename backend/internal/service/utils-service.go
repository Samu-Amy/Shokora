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
func withTransaction(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) (err error) {
	// Create transaction
	transaction, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Defer rollback (in caso di panic)
	defer func() {
		if p := recover(); p != nil {
			_ = transaction.Rollback()
			panic(p)
		} else if err != nil {
			_ = transaction.Rollback()
		}
	}()

	if err = fn(transaction); err != nil {
		return err
	}

	return transaction.Commit()
}

// - Verification Tokens -

func (service *AuthService) verifyOtp(ctx context.Context, verificationId int64, hashedOTP []byte, maxAttempts uint8, verificationType auth.VerificationType) (*store.OTPPayload, error) {
	// TODO: usare transaction (ed usare FOR UPDATE nel get?) per GetOtpData e UpdateOtpAttempts?

	// Get data
	otpQueryData, err := service.vTokensRepo.GetOtpData(ctx, verificationId, verificationType)
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
		err = service.vTokensRepo.UpdateOtpAttempts(ctx, verificationId, maxAttempts)

		// Attempts updated successfully but OTP not valid
		if err == nil {
			err = errorcodes.ErrInvalid
		}

		return nil, err
	}

	return otpQueryData, nil
}
