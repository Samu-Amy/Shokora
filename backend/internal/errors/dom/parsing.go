package domerrors

import (
	"errors"

	interrors "github.com/Samu-Amy/Shokora/internal/errors/int"
)

func ParseIntError(err error) error {
	switch {
	case errors.Is(err, interrors.IErrNotFound):
		return ErrNotFound
		// TODO: continua
	default:
		return err
	}
}
