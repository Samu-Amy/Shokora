package userservice

import (
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/store/user"
)

type UserService struct {
	userRepo user.IUserRepository
	db       *sql.DB
}

func NewService(userRepo user.IUserRepository, db *sql.DB) *UserService {
	return &UserService{userRepo, db}
}
