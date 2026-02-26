package database

import (
	"database/sql"
	"errors"

	"github.com/Samu-Amy/Shokora/internal/errorcodes"
	"github.com/lib/pq"
)

// Error codes
const (
	UNIQUE_VIOLATION_ERROR = pq.ErrorCode("23505")
	CHECK_VIOLATION_ERROR  = pq.ErrorCode("23514")
)

// Error constraints
const (
	// Users
	USERS_USER_EMAIL_UNIQUE     = "users_email_unique"
	USERS_USER_ROLE_RANGE_CHECK = "users_user_role_range_check"

	// Products
	PRODUCTS_PRICE_RANGE_CHECK    = "products_price_range_check"
	PRODUCTS_DISCOUNT_RANGE_CHECK = "products_discount_range_check"

	// Verification tokens
	VTOKENS_USER_ID_AND_VERIFICATION_TYPE_UNIQUE = "v_tokens_user_id_and_verification_type_unique"
	VTOKENS_MAGIC_LINK_TOKEN_UNIQUE              = "v_tokens_magic_link_token_hash_unique"
	VTOKENS_VERIFICATION_TYPE_RANGE_CHECK        = "v_tokens_verification_type_range_check"
	VTOKENS_OTP_ATTEMPTS_RANGE_CHECK             = "v_tokens_otp_attempts_range_check"
	VTOKENS_MAGIC_LINK_TOKEN_CHECK               = "v_tokens_magic_link_token_hash_check"

	// Refresh Tokens
	REFRESH_TOKENS_TOKEN_UNIQUE    = "refresh_tokens_token_hash_unique"
	REFRESH_TOKENS_REPLACES_UNIQUE = "refresh_tokens_replaces_unique"
)

// Check postgres error constraint
func isPostgresError(err error, errorCode pq.ErrorCode, constraint string) bool {
	if pgErr, ok := err.(*pq.Error); ok {
		return pgErr.Code == errorCode && pgErr.Constraint == constraint
	}
	return false
}

// Parse db errors into custom errorcodes when necessary or return the err
func ParseDbError(err error) error {
	switch {
	// Generic
	case errors.Is(err, sql.ErrNoRows):
		return errorcodes.ErrNotFound

	// Users
	case isPostgresError(err, UNIQUE_VIOLATION_ERROR, USERS_USER_EMAIL_UNIQUE):
		return errorcodes.ErrDuplicateEmail

	case isPostgresError(err, UNIQUE_VIOLATION_ERROR, VTOKENS_MAGIC_LINK_TOKEN_UNIQUE):
		return errorcodes.InternalErrDuplicateToken

	case isPostgresError(err, UNIQUE_VIOLATION_ERROR, REFRESH_TOKENS_REPLACES_UNIQUE):
		return errorcodes.InternalErrReusedToken // TODO: giusto?

	// TODO: implementa gestione di tutti i constraints

	default:
		return err
	}
}

// Wrap ExecContext to obtain the parsed error
func HandleExecContextResult(res sql.Result, err error) error {
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errorcodes.InternalErrNoRowsAffected
	}

	return nil
}
