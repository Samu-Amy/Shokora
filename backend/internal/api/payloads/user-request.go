package payloads

import "time"

// - Basics -

type UserDataReq struct {
	FirstName string    `json:"first_name" validate:"required,max=125"`
	LastName  string    `json:"last_name" validate:"omitempty,max=125"`       // TODO: opzionale (?)
	ImageUrl  string    `json:"image_url,omitempty" validate:"omitempty,url"` // TODO: usare url (se l'url sarà conforme al controllo)?
	BirthDate time.Time `json:"birth_date,omitempty" validate:"omitempty"`    // TODO: fare validazione
}

type EmailFieldReq struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

type PasswordFieldReq struct {
	Password string `json:"password" validate:"required,min=12,max=72"`
}

type DoublePasswordFieldReq struct {
	Password             string `json:"password" validate:"required,min=12,max=72"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=12,max=72"`
}

// - Others -

type UpdatePasswordReq struct {
	OldPassword             string `json:"old_password" validate:"required,min=12,max=72"`
	NewPassword             string `json:"new_password" validate:"required,min=12,max=72"`
	InvalidateOtherSessions bool
}
