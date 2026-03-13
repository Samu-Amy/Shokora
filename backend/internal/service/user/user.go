package userservice

import (
	"context"

	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
	"github.com/Samu-Amy/Shokora/internal/store/user"
)

// ------ CREATE USER -----

func (service *UserService) GetById(ctx context.Context, userId int64) (*user.User, error) {
	user, err := service.userRepo.GetById(ctx, userId)
	if err != nil {
		return nil, domerrors.ParseIntError(err)
	}

	return user, nil
}

// TODO: se viene usato, controllare che sia verificato
func (service *UserService) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	user, err := service.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, domerrors.ParseIntError(err)
	}

	return user, nil
}

func (service *UserService) Delete(ctx context.Context, userId int64) error {
	return domerrors.ParseIntError(service.userRepo.Delete(ctx, userId))
}
