package interrors

// - Internal Errors -

type internalErr string

func (err internalErr) Error() string {
	return string(err)
}

var (
	// Generic
	IErrNotFound       = internalErr("i_not_found")
	IErrNoRowsAffected = internalErr("i_no_rows_affected")
	IErrConflict       = internalErr("i_conflict")
	IErrExpired        = internalErr("i_expired")

	// Constraints
	IErrDuplicate                     = internalErr("i_duplicate")       // Duplicate (unique constraint)
	IErrDuplicateEmail                = internalErr("i_duplicate_email") // Duplicate email
	IErrDuplicateToken                = internalErr("i_duplicate_token") // Duplicate token
	IErrInvalid                       = internalErr("i_invalid")         // Invalid (value non valid) or range check failed)
	IErrReusedToken                   = internalErr("i_reused_token")    // Refresh Token reused (reuse detection)
	IErrTokenNotFoundOrAlreadyRevoked = internalErr("i_token_not_found_or_already_revoked")
	IErrInvalidEmailVars              = internalErr("i_invalid_email_vars") // SendEmail called with wrong variables for the template
	IErrMaxRetriesExceeded            = internalErr("i_max_retries")        // SendEmail called with wrong variables for the template
)
