package userservice

import (
	"github.com/Samu-Amy/Shokora/internal/database"
	"github.com/Samu-Amy/Shokora/internal/store/user"
	"go.uber.org/zap"
)

type UserService struct {
	txManager database.ITransactionManager
	userRepo  user.IUserRepository
	logger    *zap.SugaredLogger
}

func NewService(txManager database.ITransactionManager, userRepo user.IUserRepository, logger *zap.SugaredLogger) *UserService {
	return &UserService{txManager, userRepo, logger}
}
