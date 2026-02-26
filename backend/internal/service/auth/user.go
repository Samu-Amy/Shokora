package authservice

import (
	"context"
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/store/user"
)

/*
Executes the user registration:
  - create new user in db
    -
*/
func (service *AuthService) RegisterUser(ctx context.Context, payload payloads.RegisterUserReqPayload) (*payloads.RegisterUserResPayload, error) { // TODO: ritorna anche i dati per i cookies

	// Hash password
	hashedPassword, err := service.hashPassword(payload.Password)
	if err != nil {
		service.logger.Warnw("Error hashing password", "error", err)
		return nil, err // TODO: log errori (?)
	}

	// Create user from payload data
	user := &user.User{
		FirstName:    payload.FirstName,
		LastName:     payload.LastName,
		Email:        payload.Email,
		PasswordHash: hashedPassword,
		ImageUrl:     payload.ImageUrl,
		BirthDate:    payload.BirthDate,
	}

	// Create user in db
	if err := service.createUser(ctx, user); err != nil {
		service.logger.Warnw("Error creating user", "error", err)
		return nil, err
	}

	// Create Response Payload with user
	resPayload := &payloads.RegisterUserResPayload{
		User: *user,
	}

	// TODO: continua

	// Create Refresh Token
	refreshToken, err := service.generateNewRefreshToken(ctx, user.Id)
	if err != nil {
		service.logger.Warnw("Error generating refresh token", "error", err)
		resPayload.AuthError = true
	}

	return resPayload, nil
}

// ----- CREATE USER -----

func (service *AuthService) createUser(ctx context.Context, user *user.User) error {
	return service.txManager.WithTx(ctx, func(tx *sql.Tx) error {

		err := service.userRepo.Create(ctx, user)
		if err != nil {
			return err // TODO: fai error handling (ritorna domerrors)
		}

		// TODO: crea anche stats and settings (oppure crearle nell'update se non esistono)?

		return nil
	})

}
