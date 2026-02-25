package userservice

import (
	"context"

	"github.com/Samu-Amy/Shokora/internal/store/user"
)

// ------ CREATE USER -----

func (service *UserService) Create(ctx context.Context, user *user.User) error {
	// TODO: fare transaction per creazione user, stats and settings (oppure crearle nell'update se non esistono)?

	// TODO: fai error handling
	return service.userRepo.Create(ctx, user)
}

func (service *UserService) GetById(ctx context.Context, userId int64) (*user.User, error) {
	// TODO: fai error handling
	return service.userRepo.GetById(ctx, userId)
}

func (service *UserService) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	// TODO: fai error handling
	return service.userRepo.GetByEmail(ctx, email)
}

func (service *UserService) Delete(ctx context.Context, userId int64) error {
	// TODO: fai error handling
	return service.userRepo.Delete(ctx, userId)
}
