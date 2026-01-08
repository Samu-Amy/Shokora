package models

import "context"

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type UserRepository interface {
	Create(context.Context, *User) error
}
