package errorcodes

import "errors"

var (
	// - High Level -

	// Generic
	ErrNotFound           = errors.New("not_found")
	ErrConlflict          = errors.New("conflict") // Version conflict (on resource update)
	ErrMaxRetriesExceeded = errors.New("max_retries")

	// Auth
	ErrDuplicateEmail = errors.New("duplicate_email")

	ErrInvalid = errors.New("invalid")
	ErrExpired = errors.New("expired")

	ErrVerification = errors.New("verification_error")
	ErrEmailNotSent = errors.New("email_not_sent")

	ErrUnauthorized = errors.New("unauthorized") // User does not exists or is not verified
	ErrNotVerified  = errors.New("not_verified") // User must verify email

	// - Low Level -
	InternalErrDuplicateToken = errors.New("duplicate_token")
	ErrInvalidEmailVars       = errors.New("invalid_email_vars")

	// Old

	// ErrTokenInvalid = errors.New("token_invalid")
	// ErrTokenExpired = errors.New("token_expired")
	// ErrUserNotFound      = errors.New("user_not_found")
	// ErrUserNotAuthorized = errors.New("user_not_authorized")
	// ErrUserBlocked = errors.New("user_blocked")
)
