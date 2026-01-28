package store

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// - Errors -

var (
	ErrNotFound         = errors.New("not found")
	ErrVersionConlflict = errors.New("version conflict")
	ErrDuplicateEmail   = errors.New("duplicate email")
	ErrExpired          = errors.New("expired")
	ErrUnauthorized     = errors.New("unauthorized")
)

// - Timeouts -

const (
	medium_query_timeout = 8 * time.Second
	long_query_timeout   = 12 * time.Second
)

// - Hashing -
func HashToken(plainToken string) string {
	hash := sha256.Sum256([]byte(plainToken))
	return hex.EncodeToString(hash[:])
}

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
