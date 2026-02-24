package errorcodes

import "errors"

var (
	// - High Level -

	// Generic
	ErrNotFound           = errors.New("not_found")
	ErrConflict           = errors.New("conflict") // Version conflict (on resource update)
	ErrMaxRetriesExceeded = errors.New("max_retries")

	// Auth
	ErrDuplicateEmail = errors.New("duplicate_email")

	ErrInvalid             = errors.New("invalid")
	ErrMaxAttemptsExceeded = errors.New("max_attempts")

	ErrVerification = errors.New("verification_error")
	ErrEmailNotSent = errors.New("email_not_sent")

	ErrUnauthorized = errors.New("unauthorized") // User does not exists or is not verified
	ErrNotVerified  = errors.New("not_verified") // User must verify email

	// - Low Level -
	InternalErrExpired                       = errors.New("i_expired")
	InternalErrDuplicateToken                = errors.New("i_duplicate_token")
	InternalErrReusedToken                   = errors.New("i_reused_token")
	InternalErrTokenNotFoundOrAlreadyRevoked = errors.New("i_token_not_found_or_already_revoked")
	InternalErrInvalidEmailVars              = errors.New("i_invalid_email_vars") // SendEmail called with wrong variables for the template

	// Old

	// ErrTokenInvalid = errors.New("token_invalid")
	// ErrTokenExpired = errors.New("token_expired")
	// ErrUserNotFound      = errors.New("user_not_found")
	// ErrUserNotAuthorized = errors.New("user_not_authorized")
	// ErrUserBlocked = errors.New("user_blocked")
)
