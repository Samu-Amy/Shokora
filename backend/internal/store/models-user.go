package store

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	Id             int64  `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	HashedPassword []byte `json:"-"`
	IsVerified     bool   `json:"is_verified"`
	Role           Role   `json:"role"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// Roles
type Role uint8

const (
	RoleCustomer Role = 0
	RoleEmployee Role = 1
	RoleAdmin    Role = 2
	RoleDev      Role = 3
)

func (user *User) IsRoleValid(requiredRole Role) bool {
	return user.Role >= requiredRole
}

// Password
// type Password struct {
// 	Text *string
// 	Hash []byte
// }

// func (p *Password) Set(text string) error {
// 	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
// 	if err != nil {
// 		return err
// 	}

// 	p.Text = &text
// 	p.Hash = hash

// 	return nil
// }

// Repository
type UserRepository interface {
	// Auth main
	CreateUserAndSendVerification(context.Context, *User, string, time.Duration) error
	VerifyEmail(context.Context, string) error
	DeleteUserAndEmailVerificationToken(context.Context, int64) error

	// Auth utils
	createEmailVerification(context.Context, *sql.Tx, string, time.Duration, int64) error
	getUserFromEmailVerificationToken(context.Context, *sql.Tx, string) (*User, error)
	setUserIsVerified(context.Context, *sql.Tx, int64) error
	deleteEmailVerificationToken(context.Context, *sql.Tx, int64) error

	// Users
	Create(context.Context, *sql.Tx, *User) error
	GetById(context.Context, int64) (*User, error)
	GetByEmail(context.Context, string) (*User, error)
	Delete(context.Context, *sql.Tx, int64) error
}
