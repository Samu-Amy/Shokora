package domerrors

import "errors"

// - Domain Errors -

var (
	// Generic
	ErrNotFound           = errors.New("not_found")
	ErrConflict           = errors.New("conflict")    // Version conflict (on resource update)
	ErrMaxRetriesExceeded = errors.New("max_retries") // Sending Email, Genereting Magic Link

	// Auth
	ErrDuplicateEmail = errors.New("duplicate_email")

	ErrInvalid             = errors.New("invalid")      // Invalid data (form data, tokens)
	ErrMaxAttemptsExceeded = errors.New("max_attempts") // Too much attempts (otp)

	ErrVerification = errors.New("verification_error") //
	ErrEmailNotSent = errors.New("email_not_sent")

	ErrUnauthorized = errors.New("unauthorized") // User does not exists or is not verified
	ErrNotVerified  = errors.New("not_verified") // User must verify email
)
