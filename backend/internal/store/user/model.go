package user

import (
	"time"
)

type User struct {
	Id           int64
	FirstName    string
	LastName     string
	Email        string
	PasswordHash []byte
	ImageUrl     string
	BirthDate    time.Time
	IsVerified   bool
	IsActive     bool
	Role         Role
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
