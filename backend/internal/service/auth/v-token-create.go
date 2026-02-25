package authservice

import (
	"context"
	"errors"

	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/errorcodes"
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
