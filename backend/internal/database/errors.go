package database

import (
	"database/sql"
	"errors"

	interrors "github.com/Samu-Amy/Shokora/internal/errors/int"
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
	// usersUserEmailUnique    = "users_email_unique"          // Duplicate email
	// usersUserRoleRangeCheck = "users_user_role_range_check" // Role enum out of range

	// Verification tokens
	// vtokensUserIdAndVerificationTypeUnique = "v_tokens_user_id_and_verification_type_unique" // Unique (user_id, verification_id), usually handled by the create query
	vtokensMagicLinkTokenUnique = "v_tokens_magic_link_token_hash_unique" // Duplicate Magic Link Token
	// vtokensVerificationTypeRangeCheck      = "v_tokens_verification_type_range_check"        // Verification type enum out of range
	// vtokensOtpAttemptsRangeCheck           = "v_tokens_otp_attempts_range_check"             // Attempts ([0, 255]) out of range (smallint in uint8 range)
	// vtokensMagicLinkTokenCheck             = "v_tokens_magic_link_token_hash_check"          // Magic Link Token and expiration check (must be both null or not null)

	// Refresh Tokens
	// refreshTokensTokenUnique    = "refresh_tokens_token_hash_unique" // Duplicate Refresh Token
	refreshTokensReplacesUnique = "refresh_tokens_replaces_unique" // Duplicate "repaces" Refresh Token id (reused token)

	// Products
	// productsPriceRangeCheck    = "products_price_range_check"    // Price (>0) out of range
	// productsDiscountRangeCheck = "products_discount_range_check" // Discount ([0, 1]) out of range
)

// Check postgres error constraint
func isPostgresErrorCode(err error, errorCode pq.ErrorCode) bool {
	if pgErr, ok := err.(*pq.Error); ok {
		return pgErr.Code == errorCode
	}
	return false
}

func isPostgresError(err error, errorCode pq.ErrorCode, constraint string) bool {
	if pgErr, ok := err.(*pq.Error); ok {
		return pgErr.Code == errorCode && pgErr.Constraint == constraint
	}
	return false
}

// Parse db errors into custom interrors when necessary or return the err
func ParseDbError(err error) error {
	switch {
	// Generic
	case errors.Is(err, sql.ErrNoRows):
		return interrors.IErrNotFound

	// Verification
	case isPostgresError(err, UNIQUE_VIOLATION_ERROR, vtokensMagicLinkTokenUnique):
		return interrors.IErrDuplicateToken

	// Auth
	case isPostgresError(err, UNIQUE_VIOLATION_ERROR, refreshTokensReplacesUnique):
		return interrors.IErrReusedToken

	// Defaults
	case isPostgresErrorCode(err, CHECK_VIOLATION_ERROR):
		return interrors.IErrInvalid

	case isPostgresErrorCode(err, UNIQUE_VIOLATION_ERROR):
		return interrors.IErrDuplicate

	default:
		return err
	}
}

// Wrap ExecContext to obtain the parsed error
func HandleExecContextResult(res sql.Result, err error) error {
	if err != nil {
		return ParseDbError(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return interrors.IErrNoRowsAffected
	}

	return nil
}
