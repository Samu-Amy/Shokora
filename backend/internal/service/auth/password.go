package authservice

import "golang.org/x/crypto/bcrypt"

func (service *AuthService) HashPassword(plainPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), service.config.PasswordHashingCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}
