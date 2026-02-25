package authservice

import (
	"context"

	"github.com/Samu-Amy/Shokora/internal/store/user"
)

// ------ CREATE USER -----

func (service *AuthService) CreateUser(ctx context.Context, user *user.User) error {
	// TODO: fare transaction per creazione user, stats and settings (oppure crearle nell'update se non esistono)?

	return service.userRepo.Create(ctx, user)
}
