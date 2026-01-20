package payload

/*
Password validation:
- min 8 chars
- max 72 chars
*/

type RegisterUserPayload struct {
	FirstName string `json:"first_name" validate:"required,max=125"`
	LastName  string `json:"last_name" validate:"required,max=125"`
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=8,max=72"` // TODO: aggiungere altri controlli?
	// BirthDate string `json:"birth_date,omitempty" validate:"omitempty"` // TODO: togliere omitempty (?)
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}
