package userservice

import (
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/store/user"
)

type UserService struct {
	userRepo user.IUserRepository
	db       *sql.DB
}

func NewService(db *sql.DB, userRepo user.IUserRepository) *UserService {
	return &UserService{userRepo, db}
}
