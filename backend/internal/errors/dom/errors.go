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
	ErrInternalError      = domainErr("internal_error")
	ErrBadParam           = domainErr("bad_param")

	// Auth
	ErrDuplicateEmail      = domainErr("duplicate_email")
	ErrInvalid             = domainErr("invalid")          // Invalid data (form data, tokens)
	ErrMaxAttemptsExceeded = domainErr("max_attempts")     // Too much attempts (otp)
	ErrUnauthorized        = domainErr("unauthorized")     // User does not exists or is not verified
	ErrForbidden           = domainErr("forbidden")        // User does not have the necessary permissions
	ErrNotVerified         = domainErr("not_verified")     // User must verify email
	ErrOnlyGoogleAuth      = domainErr("only_google_auth") // User can't use email and password, but only Google OAuth

	// Validation
	ErrInvalidDate = domainErr("invalid_date") // The date (e.g. birthdate) is invalid
)

// CAUTION: does not work with wrapping (fmt.Errorf("...%w...", err))
func IsDomainErr(err error) bool {
	_, ok := err.(domainErr)
	return ok
}
