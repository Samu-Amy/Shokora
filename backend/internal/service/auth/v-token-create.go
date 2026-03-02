package authservice

import (
	"context"
	"errors"

	"github.com/Samu-Amy/Shokora/internal/auth"
	interrors "github.com/Samu-Amy/Shokora/internal/errors/int"
)

// ----- CREATE VERIFICATION TOKENS -----

func (service *AuthService) createVerificationTokensWithRetries(ctx context.Context, userId int64, verificationType auth.VerificationType) (*auth.VerificationTokens, error) {

	// Generate tokens
	verificationTokens, err := service.tokenAuthenticator.CreateVerificationTokens(verificationType)
	if err != nil {
		service.logger.Warnw("Error generating verification tokens", "error", err)
		return nil, err
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, auth.RegenerateTokenTimeout)
	defer cancel()

	// Create tokens in db
	err = service.vTokenRepo.Create(ctxWithTimeout, userId, verificationTokens)
	if err == nil {
		return verificationTokens, nil // Tokens Created (OK)

	} else if !errors.Is(err, interrors.IErrDuplicateToken) {
		service.logger.Warnw("Error creating verification tokens in db", "error", err)
		return nil, err // Error (can't retry)
	}

	// Retries (it's MaxRetries - 1 because one try is already done)
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
		case errors.Is(err, interrors.IErrDuplicateToken):
			err = service.tokenAuthenticator.RegenerateMagicLinkToken(verificationTokens)
			if err != nil {
				service.logger.Warnw("Error rigenerating verification tokens", "error", err)
				continue // skip iteration
			}

			err = service.vTokenRepo.Create(ctx, userId, verificationTokens)
			if err == nil {
				return verificationTokens, nil // Tokens Created (OK)
			}

		default:
			service.logger.Warnw("Error during verification tokens creation retries", "error", err)
			return nil, err // Error is not solvable (not "duplicate token") -> return it
		}
	}

	return nil, interrors.IErrMaxRetriesExceeded // Couldn't regenerate and save token successfully -> return error "max retries exceeded"
}
