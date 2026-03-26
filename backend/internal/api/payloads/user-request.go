package payloads

// - Basics -

// TODO: finisci validazione (fai anche per stringhe generiche, tipo nomi prodotti e descrizioni), togli validazione da service

type UserDataReq struct {
	FirstName string `json:"first_name" validate:"required,min=1,max=125,valid-name"`
	LastName  string `json:"last_name" validate:"omitempty,max=125,valid-name"`
	Birthday  string `json:"birthday,omitempty" validate:"omitempty,valid-birthday"`
	// ImageUrl  string    `json:"image_url,omitempty" validate:"omitempty,url,valid-chars"` // usare url (se l'url sarà conforme al controllo)?
}

type EmailFieldReq struct {
	Email string `json:"email" validate:"required,email,max=255"` // TODO: migliora controllo sull'email
}

type PasswordFieldReq struct {
	Password string `json:"password" validate:"required,min=12,max=72,no-edge-spaces,valid-chars,no-common-password"` // bcrypt clamp to 72 chars
}

type DoublePasswordFieldReq struct {
	Password             string `json:"password" validate:"required,min=12,max=72,no-edge-spaces,valid-chars,no-common-password"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=12,max=72,no-edge-spaces,valid-chars,no-common-password,eqfield=Password"`
}

// - Others -

type UpdatePasswordReq struct {
	OldPassword             string `json:"old_password" validate:"required,min=12,max=72,no-edge-spaces,valid-chars,no-common-password"`
	NewPassword             string `json:"new_password" validate:"required,min=12,max=72,no-edge-spaces,valid-chars,no-common-password,nefield=OldPassword"`
	InvalidateOtherSessions bool
}
