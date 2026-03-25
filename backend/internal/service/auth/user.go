package authservice

import (
	"context"
	"database/sql"
	"errors"
	"time"

	interrors "github.com/Samu-Amy/Shokora/internal/errors/int"
	"github.com/Samu-Amy/Shokora/internal/store/user"
	"golang.org/x/crypto/bcrypt"
)

// - Create -

// Create user and related tables (e.g. settings, stats, achievements, copons)
func (service *AuthService) createUser(ctx context.Context, user *user.User) error {
	return service.txManager.WithTx(ctx, func(tx *sql.Tx) error { // TODO: usare transaction oppure creare solo user e creare le righe nelle altre tabelle a parte (e se falliscono si creano quando vengono usate (però non si possono ottenere))

		err := service.userRepo.Create(ctx, tx, user)
		if err != nil {
			service.logger.Warnw("Error creating user in db", "error", err)
			return err
		}

		_, err = service.userSettingsRepo.Create(ctx, tx, user.Id)
		if err != nil {
			service.logger.Warnw("Error creating user settings", "error", err)
			return err
		}

		// TODO: crea anche stats (oppure crearle nell'update se non esistono)?

		return nil
	})

}

// - Get -

func (service *AuthService) getUser(ctx context.Context, email string, plainPassword string) (*user.User, error) {

	// Get user from db
	user, err := service.userRepo.GetByEmail(ctx, email)
	if err != nil {
		service.logger.Warnw("Error getting user from db", "error", err)
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

		service.logger.Warnw("Password compare error", "error", err)
		return nil, err
	}

	// Check user
	if !user.IsActive {
		return nil, interrors.IErrUnauthorized
	}

	if !user.IsVerified {
		return nil, interrors.IErrNotVerified
	}

	// if user.HasTwoFactorAuth { // TODO: aggiungi two factor auth check (magari join con settings table?)
	// return nil, interrors.IErrTwoFactorAuthReqired
	// }

	return user, nil
}

// ----- UTILS -----

// Birthday conversion
func convertBirthdayToTime(birthdayStr string) (time.Time, error) {
	// The layout is "02" -> day, "01" -> month, "2006" -> year, using 2000 because is a leap year (to avoid errors in the case of February 29)
	return time.Parse("02-01-2006", birthdayStr+"-2000")
}
