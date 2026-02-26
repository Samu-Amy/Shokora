package database

import (
	"net/http"
	"strconv"
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
