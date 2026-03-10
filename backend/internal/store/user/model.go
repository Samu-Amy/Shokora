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
	Permissions  Permission
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

// - Permissions -
type Permission uint32

const (
	ProductAdd Permission = 1 << iota
	ProductUpdate
	ProductDelete
	StockProductAdd
	StockProductUpdate
	StockProductDelete
	OrderUpdate
	OrderDelete
)

func (user *User) AddPermission(permission Permission) {
	user.Permissions |= permission
}

func (user *User) HasPermission(permission Permission) bool {
	return user.Permissions&permission != 0
}

func (user *User) RemovePermission(permission Permission) {
	user.Permissions &^= permission
}
