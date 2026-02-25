package userservice

import (
	"database/sql"

	"github.com/Samu-Amy/Shokora/internal/store/user"
)

type UserService struct {
	userRepo user.UserRepositoryI
	db       *sql.DB
}

func NewService(userRepo user.UserRepositoryI, db *sql.DB) *UserService {
	return &UserService{userRepo, db}
}
