package domerrors

import (
	"errors"

	interrors "github.com/Samu-Amy/Shokora/internal/errors/int"
)

func ParseIntError(err error) error {
	switch {
	case errors.Is(err, interrors.IErrNotFound):
		return ErrNotFound

	case errors.Is(err, interrors.IErrInvalid):
		return ErrInvalid

	case errors.Is(err, interrors.IErrNoRowsAffected):
		return ErrInvalid // TODO: cambia?

	case errors.Is(err, interrors.IErrDuplicateToken):
		return ErrInvalid // TODO: cambia (dipende, può essere not valid (se token fornito dall'utente) o internal error (se token generato internamente))

	case errors.Is(err, interrors.IErrExpired):
		return ErrInvalid

	case errors.Is(err, interrors.IErrDuplicate):
		return ErrInvalid // TODO: cambia (dipende, può essere not valid o internal error)

	case errors.Is(err, interrors.IErrConflict):
		return ErrConflict

	case errors.Is(err, interrors.IErrDuplicateEmail):
		return ErrDuplicateEmail

	case errors.Is(err, interrors.IErrReusedToken):
		return ErrUnauthorized // TODO: cambia?

	case errors.Is(err, interrors.IErrTokenNotFoundOrAlreadyRevoked):
		return ErrUnauthorized // TODO: cambia?

	case errors.Is(err, interrors.IErrInvalidEmailVars):
		return ErrInternalError

	case errors.Is(err, interrors.IErrMaxRetriesExceeded):
		return ErrMaxRetriesExceeded

	case errors.Is(err, interrors.IErrMaxAttemptsExceeded):
		return ErrMaxAttemptsExceeded

	case errors.Is(err, interrors.IErrNotVerified):
		return ErrNotVerified

	default:
		return ErrInternalError
	}
}
