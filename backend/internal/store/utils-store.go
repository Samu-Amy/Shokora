package store

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Samu-Amy/Shokora/internal/errorcodes"
	"github.com/lib/pq"
)

// - Errors -

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
func parseDbError(err error) error {
	switch {
	// Generic
	case errors.Is(err, sql.ErrNoRows):
		return errorcodes.ErrNotFound

	// Users
	case isPostgresError(err, UNIQUE_VIOLATION_ERROR, USERS_USER_EMAIL_UNIQUE):
		return errorcodes.ErrDuplicateEmail

	case isPostgresError(err, UNIQUE_VIOLATION_ERROR, VTOKENS_MAGIC_LINK_TOKEN_UNIQUE):
		return errorcodes.InternalErrDuplicateToken

	default:
		return err
	}
}

// Wrap ExecContext to obtain the parsed error
func handleExecContextResult(res sql.Result, err error) error {
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

// - Timeouts -

const (
	MEDIUM_QUERY_TIMEOUT = 8 * time.Second
	LONG_QUERY_TIMEOUT   = 12 * time.Second
)

// - Pagination -

type QueryPaginationOptions struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=50"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
}

func (queryPaginationOptions QueryPaginationOptions) Parse(r *http.Request) (QueryPaginationOptions, error) {
	queryString := r.URL.Query()

	// Limit
	limit := queryString.Get("limit")
	if limit != "" {
		lim, err := strconv.Atoi(limit)
		if err != nil {
			return queryPaginationOptions, err
		}

		queryPaginationOptions.Limit = lim
	}

	// Offset
	offset := queryString.Get("offset")
	if offset != "" {
		off, err := strconv.Atoi(offset)
		if err != nil {
			return queryPaginationOptions, err
		}

		queryPaginationOptions.Offset = off
	}

	// Sort
	sort := queryString.Get("sort")
	if sort != "" {
		queryPaginationOptions.Sort = sort
	}

	return queryPaginationOptions, nil
}

// - Filters -

type MenuFilters struct {
	Search string   `json:"search" validate:"max=100"`
	Badges []string `json:"badges"` // TODO: usare enum e come fare validate (?)
}

func (menuFilters MenuFilters) Parse(r *http.Request) (MenuFilters, error) {
	queryString := r.URL.Query()

	// Search
	search := queryString.Get("search")
	if search != "" {
		menuFilters.Search = search
	}

	// Badges
	badges := queryString.Get("badges")
	if badges != "" {
		menuFilters.Badges = strings.Split(badges, ",")
	}

	return menuFilters, nil
}

// type ShopFilters struct {
// 	Search string   `json:"search" validate:"max=100"`
// 	Badges []string `json:"badges"` // TODO: usare enum e come fare validate (?)
// }

// func (shopFilters ShopFilters) Parse(r *http.Request) (ShopFilters, error) {
// 	queryString := r.URL.Query()

// 	// Search
// 	search := queryString.Get("search")
// 	if search != "" {
// 		shopFilters.Search = search
// 	}

// 	// Badges
// 	badges := queryString.Get("badges")
// 	if badges != "" {
// 		shopFilters.Badges = strings.Split(badges, ",")
// 	}

// 	return shopFilters, nil
// }

type ProductsFilters struct {
	Search string   `json:"search" validate:"max=100"`
	Badges []string `json:"badges"` // TODO: usare enum e come fare validate (?)
	Since  string   `json:"since"`
	Until  string   `json:"until"`
}

func (productsFilters ProductsFilters) Parse(r *http.Request) (ProductsFilters, error) {
	queryString := r.URL.Query()

	// Search
	search := queryString.Get("search")
	if search != "" {
		productsFilters.Search = search
	}

	// Badges
	badges := queryString.Get("badges")
	if badges != "" {
		productsFilters.Badges = strings.Split(badges, ",")
	}

	// Since
	since := queryString.Get("since")
	if since != "" {
		productsFilters.Since = parseTime(since)
	}

	// Until
	until := queryString.Get("until")
	if until != "" {
		productsFilters.Until = parseTime(until)
	}

	return productsFilters, nil
}

func parseTime(strTime string) string {
	t, err := time.Parse(time.DateTime, strTime)
	if err != nil {
		return ""
	}
	return t.Format(time.DateTime)
}
