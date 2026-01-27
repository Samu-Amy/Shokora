package payload

import "time"

/*
Password validation:
- min 8 chars
- max 72 chars
*/

type RegisterUserPayload struct {
	FirstName string    `json:"first_name" validate:"required,max=125"`
	LastName  string    `json:"last_name" validate:"required,max=125"`
	Email     string    `json:"email" validate:"required,email,max=255"`
	Password  string    `json:"password" validate:"required,min=8,max=72"`    // TODO: aggiungere altri controlli?
	ImageUrl  string    `json:"image_url,omitempty" validate:"omitempty,url"` // TODO: usare url (se l'url sarà conforme al controllo)?
	BirthDate time.Time `json:"birth_date,omitempty" validate:"omitempty"`    // TODO: fare validazione
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// TODO: fare validazione custom (tipo quella sotto)?
// validate := validator.New()

// validate.RegisterValidation("birthdate", func(fl validator.FieldLevel) bool {
//     date, ok := fl.Field().Interface().(time.Time)
//     if !ok || date.IsZero() {
//         return true // omitempty
//     }

//     // Non può essere nel futuro
//     if date.After(time.Now()) {
//         return false
//     }

//     // Età minima 13 anni (esempio)
//     minAge := time.Now().AddDate(-13, 0, 0)
//     return date.Before(minAge)
// })
