package authservice

import (
	"context"
	"errors"

	"github.com/Samu-Amy/Shokora/internal/auth"
	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
	interrors "github.com/Samu-Amy/Shokora/internal/errors/int"
	"github.com/Samu-Amy/Shokora/internal/store/user"
)

// ----- CREATE VERIFICATION TOKENS -----

func (service *AuthService) createVerificationTokensWithRetries(ctx context.Context, user *user.User) (*auth.VerificationTokens, error) {

	// Generate tokens
	verificationTokens, err := service.tokenAuthenticator.CreateVerificationTokens(auth.EmailVerification)
	if err != nil {
		service.logger.Warnw("Error generating verification tokens", "error", err)
		return nil, err
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, auth.RegenerateTokenTimeout)
	defer cancel()

	// Create tokens in db
	verificationId, err := service.vTokenRepo.CreateTokens(ctxWithTimeout, user.Id, verificationTokens) // TODO: controlla che le scadenze siano giuste
	if err == nil {
		return verificationTokens, nil // OK, return no error

	} else if !errors.Is(err, interrors.IErrDuplicateToken) {
		service.logger.Warnw("Error creating verification tokens in db", "error", err)
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
		case errors.Is(err, interrors.IErrDuplicateToken):
			err = service.tokenAuthenticator.RegenerateMagicLinkToken(verificationTokens)
			if err != nil {
				service.logger.Warnw("Error rigenerating verification tokens", "error", err)
				continue // skip iteration
			}

			verificationId, err = service.vTokenRepo.CreateTokens(ctx, user.Id, verificationTokens)
			if err == nil {
				service.logger.Warnw("Error ricreating verification tokens in db", "error", err)
				return verificationTokens, nil // OK, return no error
			}

		default:
			service.logger.Warnw("Error in create verification tokens retries", "error", err)
			return nil, err // Error is not solvable (not "duplicate token") -> return it
		}
	}

	return nil, domerrors.ErrMaxRetriesExceeded // Couldn't regenerate and save token successfully -> return error "max retries exceeded"
}
