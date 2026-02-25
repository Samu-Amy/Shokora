package user

import (
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

// - Role -

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
