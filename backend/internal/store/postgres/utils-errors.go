package postgres

import "errors"

var (
	ErrNotFound         = errors.New("resource not found")
	ErrVersionConlflict = errors.New("version conflict")
)
