package interrors

// - Internal Errors -

type internalErr string

func (err internalErr) Error() string {
	return string(err)
}

var (
	IErrNoRowsAffected                = internalErr("i_no_rows_affected")
	IErrExpired                       = internalErr("i_expired")
	IErrDuplicateToken                = internalErr("i_duplicate_token")
	IErrReusedToken                   = internalErr("i_reused_token")
	IErrTokenNotFoundOrAlreadyRevoked = internalErr("i_token_not_found_or_already_revoked")
	IErrInvalidEmailVars              = internalErr("i_invalid_email_vars") // SendEmail called with wrong variables for the template
	IErrMaxRetriesExceeded            = internalErr("i_max_retries")        // SendEmail called with wrong variables for the template
)
