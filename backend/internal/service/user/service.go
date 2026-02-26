package userservice

import (
	"github.com/Samu-Amy/Shokora/internal/database"
	"github.com/Samu-Amy/Shokora/internal/store/user"
)

type UserService struct {
	txManager database.ITransactionManager
	userRepo  user.IUserRepository
}

func NewService(txManager database.ITransactionManager, userRepo user.IUserRepository) *UserService {
	return &UserService{txManager, userRepo}
}
