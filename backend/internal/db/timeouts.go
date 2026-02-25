package db

import "time"

const (
	MEDIUM_QUERY_TIMEOUT = 8 * time.Second
	LONG_QUERY_TIMEOUT   = 12 * time.Second
)
