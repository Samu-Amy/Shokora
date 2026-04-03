package authservice

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
	interrors "github.com/Samu-Amy/Shokora/internal/errors/int"
	"github.com/Samu-Amy/Shokora/internal/store/user"
	"golang.org/x/crypto/bcrypt"
)

// - Create -

// Create user and related tables (e.g. settings, stats, achievements, coupons)
func (service *AuthService) createUser(ctx context.Context, tx *sql.Tx, user *user.User) error {
	err := service.userRepo.Create(ctx, tx, user)
	if err != nil {
		service.logger.Warnw("Error creating user in db", "error", err)
		if errors.Is(err, interrors.IErrDuplicate) {
			return interrors.IErrDuplicateEmail
		}
		return err
	}

	_, err = service.userSettingsRepo.Create(ctx, tx, user.Id)
	if err != nil {
		service.logger.Warnw("Error creating user settings", "error", err)
		return err
	}

	// TODO: crea anche stats (oppure crearle nell'update se non esistono)?

	return nil
}

// - Get -

func (service *AuthService) getUser(ctx context.Context, email string, plainPassword string) (*user.User, error) {

	email = strings.TrimSpace(email)
	plainPassword = strings.TrimSpace(plainPassword)

	// Get user from db
	user, err := service.userRepo.GetByEmail(ctx, email)
	if err != nil {
		service.logger.Warnw("Error getting user from db ", "error", err)

		if errors.Is(err, interrors.IErrNotFound) {
			return nil, interrors.IErrInvalid
		}

		return nil, err
	}

	// Check password
	if err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(plainPassword)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) { // TODO: aggiungere numero massimo di tentativi (?)
			return nil, interrors.IErrInvalid
		}

		if errors.Is(err, bcrypt.ErrHashTooShort) && user.GoogleId != nil {
			service.logger.Warnw("Only google auth is valid ", "error", err)
			return nil, domerrors.ErrOnlyGoogleAuth
		}

		service.logger.Warnw("Password compare error ", "error", err)
		return nil, err
	}

	// Check user
	if !user.IsActive {
		service.logger.Warn("User not active")
		return nil, interrors.IErrUnauthorized
	}

	// TODO: aggiungi check per blocked (se usato)

	if !user.IsVerified {
		service.logger.Warn("User not verified")
		return user, interrors.IErrNotVerified
	}

	// Check 2FA
	hasTwoFactorAuth, err := service.userSettingsRepo.GetHasTwoFactorAuthById(ctx, user.Id)
	if err != nil {
		service.logger.Warnw("Error getting hasTwoFactorAuth", "error", err)
		return nil, err
	}

	if hasTwoFactorAuth {
		return user, interrors.IErrTwoFactorAuthReqired
	}

	return user, nil
}

// ----- UTILS -----

// Birthday conversion
func convertBirthdayToTime(birthdayStr string) (time.Time, error) {
	// The layout is "02" -> day, "01" -> month, "2006" -> year, using 2000 because is a leap year (to avoid errors in the case of February 29)
	return time.Parse("02-01-2006", birthdayStr+"-2000")
}
