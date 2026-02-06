package store

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	Id           int64     `json:"id"` // Generated
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"-"`
	ImageUrl     string    `json:"image_url"`
	BirthDate    time.Time `json:"birth_date"`
	IsVerified   bool      `json:"is_verified"` // Default false
	IsActive     bool      `json:"is_active"`   // Default true
	Role         Role      `json:"role"`        // Default 0
	CreatedAt    time.Time `json:"created_at"`  // Default now()
	UpdatedAt    time.Time `json:"updated_at"`  // Default now()
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

// TODO: sistema (qua solo metodi "strettamente legati" a users (tabella) che fanno le query, poi quelli "composti" o con logica (es. retry) li si crea nel service usando questi)

// Repository
type UserRepositoryI interface {
	// Auth main
	// ResendEmailVerificationEmail(ctx context.Context, email string) error // TODO: cambia nome
	// DeleteUserAndEmailVerificationToken(ctx context.Context, userId int64) error

	// Auth utils
	getUserFromEmailVerificationToken(ctx context.Context, transaction *sql.Tx, plainToken string) (*User, error)
	setUserIsVerified(ctx context.Context, transaction *sql.Tx, userId int64) error
	deleteEmailVerificationToken(ctx context.Context, transaction *sql.Tx, userId int64) error

	// Users
	Create(ctx context.Context, user *User) error
	GetById(ctx context.Context, userId int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Delete(ctx context.Context, transaction *sql.Tx, userId int64) error
}
