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
	// ImageUrl     string
	Birthday time.Time // Only day and month
	// ParentalConsent bool // TODO: implementa per i minori di 14 (solo italia) o 16 anni (GDPR) (i genitori devono verificare tramite la loro email (va modificato il funzionamento della verifica email, però i genitori devono quindi registrarsi?) e devono poter revocare il consenso, quindi può servire una tabella che lega id genitori a id figli per poter avere accesso ai loro dati (e poter collegare id di uno o più genitoria uno o più figli (dalle impostazioni?) - però una volta che superano i 16 anni dovrebbero poter "scollegare" l'account da quello dei genitori)
	IsVerified  bool
	IsActive    bool // TODO: usare IsActive per abilitare/disabilitare account (l'utente e/o genitrore dell'utente) e fare IsBlocked per bloccare l'account (gestito da admin)
	Role        Role
	Permissions Permission
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
	// Products
	EmplPermProductAdd Permission = 1 << iota
	EmplPermProductUpdate
	EmplPermProductDelete

	// Stock Products
	EmplPermStockProductAdd
	EmplPermStockProductUpdate
	EmplPermStockProductDelete

	// Orders
	EmplPermOrderGet
	EmplPermOrderUpdate
	EmplPermOrderDelete
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
