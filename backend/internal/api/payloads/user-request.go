package payloads

import "time"

// - Basics -

// TODO: finisci validazione (fai anche per stringhe generiche, tipo nomi prodotti e descrizioni), togli validazione da service

type UserDataReq struct {
	FirstName string    `json:"first_name" validate:"required,max=125,valid-name"`
	LastName  string    `json:"last_name" validate:"omitempty,max=125,valid-name"`
	Birthday  time.Time `json:"birthday,omitempty" validate:"omitempty"`
	// ImageUrl  string    `json:"image_url,omitempty" validate:"omitempty,url"` // usare url (se l'url sarà conforme al controllo)?
}

type EmailFieldReq struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

type PasswordFieldReq struct {
	Password string `json:"password" validate:"required,min=12,max=72,no-edge-spaces,valid-password,no-common-password"`
}

type DoublePasswordFieldReq struct {
	Password             string `json:"password" validate:"required,min=12,max=72,no-edge-spaces,valid-password,no-common-password"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=12,max=72,no-edge-spaces,valid-password,no-common-password,eqfield=Password"`
}

// - Others -

type UpdatePasswordReq struct {
	OldPassword             string `json:"old_password" validate:"required,min=12,max=72,no-edge-spaces,valid-password,no-common-password"`
	NewPassword             string `json:"new_password" validate:"required,min=12,max=72,no-edge-spaces,valid-password,no-common-password,nefield=OldPassword"`
	InvalidateOtherSessions bool
}
