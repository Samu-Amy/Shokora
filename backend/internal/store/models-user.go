package store

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         int64    `json:"id"`
	FirstName  string   `json:"first_name"`
	LastName   string   `json:"last_name"`
	Email      string   `json:"email"`
	Password   password `json:"-"`
	IsVerified bool     `json:"is_verified"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash

	return nil
}

type UserRepository interface {
	// Auth main
	CreateUserAndSendVerification(context.Context, *User, string, time.Duration) error
	VerifyEmail(context.Context, string) error

	// Auth utils
	createEmailVerification(context.Context, *sql.Tx, string, time.Duration, int64) error
	getUserFromEmailVerificationToken(context.Context, *sql.Tx, string) (*User, error)
	setUserIsVerified(context.Context, *sql.Tx, int64) error
	deleteEmailVerificationToken(context.Context, *sql.Tx, int64) error

	// Users
	Create(context.Context, *sql.Tx, *User) error
	GetById(context.Context, int64) (*User, error)
}
