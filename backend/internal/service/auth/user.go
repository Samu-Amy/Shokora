package authservice

import (
	"context"
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/store/user"
)

// Create user and related tables (e.g. settings, stats, achievements, copons)
func (service *AuthService) createUser(ctx context.Context, user *user.User) error {
	return service.txManager.WithTx(ctx, func(tx *sql.Tx) error { // TODO: usare transaction oppure creare solo user e creare le righe nelle altre tabelle a parte (e se falliscono si creano quando vengono usate (però non si possono ottenere))

		err := service.userRepo.Create(ctx, user)
		if err != nil {
			service.logger.Warnw("Error creating user in db", "error", err)
			return err
		}

		// TODO: crea anche stats and settings (oppure crearle nell'update se non esistono)?

		return nil
	})

}
