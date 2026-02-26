package interrors

import "errors"

// - Internal Errors -

var (
	IErrNoRowsAffected                = errors.New("i_no_rows_affected")
	IErrExpired                       = errors.New("i_expired")
	IErrDuplicateToken                = errors.New("i_duplicate_token")
	IErrReusedToken                   = errors.New("i_reused_token")
	IErrTokenNotFoundOrAlreadyRevoked = errors.New("i_token_not_found_or_already_revoked")
	IErrInvalidEmailVars              = errors.New("i_invalid_email_vars") // SendEmail called with wrong variables for the template
)
