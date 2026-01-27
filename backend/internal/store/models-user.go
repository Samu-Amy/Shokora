package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/Samu-Amy/Shokora/internal/auth"
)

type User struct {
	Id           int64     `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"-"`
	ImageUrl     string    `json:"image_url"`
	BirthDate    time.Time `json:"birth_date"`
	IsVerified   bool      `json:"is_verified"`
	IsActive     bool      `json:"is_active"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
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

// Repository
type UserRepository interface {
	// Auth main
	CreateUserAndSendVerification(ctx context.Context, user *User, verificationTokens *auth.VerificationTokens) error
	ResendEmailVerificationEmail(ctx context.Context, email string) error
	VerifyEmail(ctx context.Context, plainToken string) error
	DeleteUserAndEmailVerificationToken(ctx context.Context, userId int64) error

	// Auth utils
	createEmailVerification(ctx context.Context, transaction *sql.Tx, verificationTokens *auth.VerificationTokens, userId int64) error
	getUserFromEmailVerificationToken(ctx context.Context, transaction *sql.Tx, plainToken string) (*User, error)
	setUserIsVerified(ctx context.Context, transaction *sql.Tx, userId int64) error
	deleteEmailVerificationToken(ctx context.Context, transaction *sql.Tx, userId int64) error

	// Users
	Create(ctx context.Context, transaction *sql.Tx, user *User) error
	GetById(ctx context.Context, userId int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Delete(ctx context.Context, transaction *sql.Tx, userId int64) error
}
