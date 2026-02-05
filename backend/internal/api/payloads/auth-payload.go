package payloads

import (
	"time"

	"github.com/Samu-Amy/Shokora/internal/store"
)

// - User Payload Fields -

type UserData struct {
	FirstName string    `json:"first_name" validate:"required,max=125"`
	LastName  string    `json:"last_name" validate:"omitempty,max=125"`       // TODO: opzionale (?)
	ImageUrl  string    `json:"image_url,omitempty" validate:"omitempty,url"` // TODO: usare url (se l'url sarà conforme al controllo)?
	BirthDate time.Time `json:"birth_date,omitempty" validate:"omitempty"`    // TODO: fare validazione
}

type EmailField struct {
	Email string `json:"email" validate:"required,email,max=255"`
}

type PasswordField struct {
	Password string `json:"password" validate:"required,min=8,max=72"` // TODO: aggiungere altri controlli?
}

// - Verification Payload Fields -

type VerificationIdField struct {
	VerificationId int64 `json:"verification_id" validate:"required,gte=0"`
}

type OTPField struct {
	OTP string `json:"otp" validate:"required,min=4,max=10"`
}

// - Request Payloads -

type RegisterUserReqPayload struct {
	UserData
	EmailField
	PasswordField
}

type OTPVerificationReqPayload struct {
	VerificationIdField
	OTPField
}

type LoginUserReqPayload struct {
	EmailField
	PasswordField
}

// - Response Payloads -

type UserResPayload struct {
	Id         int64      `json:"id"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Email      string     `json:"email"`
	ImageUrl   string     `json:"image_url"`
	BirthDate  time.Time  `json:"birth_date"`
	IsVerified bool       `json:"is_verified"`
	Role       store.Role `json:"role"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type RegisterUserResPayload struct {
	User           UserResPayload
	VerificationId *int64 `json:"verification_id,omitempty"`
	Error          string `json:"error,omitempty"`
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
