package domerrors

// - Domain Errors -

type domainErr string

func (err domainErr) Error() string {
	return string(err)
}

var (
	// Generic
	ErrNotFound           = domainErr("not_found")
	ErrConflict           = domainErr("conflict")    // Version conflict (on resource update)
	ErrMaxRetriesExceeded = domainErr("max_retries") // Sending Email, Genereting Magic Link

	// Auth
	ErrDuplicateEmail      = domainErr("duplicate_email")
	ErrInvalid             = domainErr("invalid")      // Invalid data (form data, tokens)
	ErrMaxAttemptsExceeded = domainErr("max_attempts") // Too much attempts (otp)
	ErrVerification        = domainErr("verification_error")
	ErrEmailNotSent        = domainErr("email_not_sent")
	ErrUnauthorized        = domainErr("unauthorized") // User does not exists or is not verified
	ErrForbidden           = domainErr("forbidden")    // User does not have the necessary permissions
	ErrNotVerified         = domainErr("not_verified") // User must verify email
)

// CAUTION: does not work with wrapping (fmt.Errorf("...%w...", err))
func IsDomainErr(err error) bool {
	_, ok := err.(domainErr)
	return ok
}
