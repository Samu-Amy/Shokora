package authservice

import (
	"context"
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/store/user"
)

func (service *AuthService) RegisterUser(ctx context.Context, payload payloads.RegisterUserReqPayload) (*payloads.RegisterUserResPayload, error) {
	// Hash password
	hashedPassword, err := service.hashPassword(payload.Password)
	if err != nil {
		return nil, err
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

	// Create Response Payload
	resPayload := &payloads.RegisterUserResPayload{}

	// Create user in db
	if err := service.createUser(ctx, user); err != nil {
		return nil, err
	}

	resPayload.User = payloads.CreateUserResPayload(user) // Add user to payload

	// Create Refresh Token
	refreshToken, err := app.generateNewRefreshToken(ctx, user.Id)
	if err != nil {
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
