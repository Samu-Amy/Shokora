package authservice

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
	interrors "github.com/Samu-Amy/Shokora/internal/errors/int"
	"golang.org/x/crypto/bcrypt"
)

// - Update -

func (service *AuthService) UpdatePassword(ctx context.Context, userId, sessionId int64, payload *payloads.UpdatePasswordReq) error {

	err := service.txManager.WithTx(ctx, func(tx *sql.Tx) error {

		// Get old password from userId
		hashedOldPassword, err := service.userRepo.GetPasswordForUpdate(ctx, tx, userId)
		if err != nil {
			return err
		}

		// Check old password
		if err = bcrypt.CompareHashAndPassword(hashedOldPassword, []byte(strings.TrimSpace(payload.OldPassword))); err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				return interrors.IErrInvalid
			}

			service.logger.Warnw("Password compare error", "userId", userId, "error", err)
			return err
		}

		// Hash new password
		hashedNewPassword, err := service.hashPassword(strings.TrimSpace(payload.NewPassword))
		if err != nil {
			return err
		}

		// Update password
		err = service.userRepo.UpdatePassword(ctx, tx, userId, hashedNewPassword)
		if err != nil {
			return err
		}

		// Invalidate sessions (delete all session with this userId)
		if payload.InvalidateOtherSessions {
			if sessionId == 0 {
				return interrors.IErrInvalid
			}

			err = service.userSessionRepo.DeleteOtherUserSessions(ctx, tx, userId, sessionId)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return domerrors.ParseIntError(err)
}
