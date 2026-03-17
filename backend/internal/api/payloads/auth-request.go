package payloads

import "github.com/google/uuid"

// Auth
type RegisterUserReq struct {
	UserDataReq
	EmailFieldReq
	PasswordFieldReq
}

type LoginUserReq struct {
	EmailFieldReq
	PasswordFieldReq
}

// Verification
type OTPVerificationReq struct {
	VerificationId uuid.UUID `json:"verification_id" validate:"gte=0"`
	OTP            string    `json:"otp" validate:"required,min=4,max=10"`
}

type PasswordResetReq struct {
	EmailFieldReq
}

// TODO: fare validazione custom (tipo quella sotto)?

// validate := validator.New()

// validate.RegisterValidation("birthdate", func(fl validator.FieldLevel) bool {
//     date, ok := fl.Field().Interface().(time.Time)
//     if !ok || date.IsZero() {
//         return true // omitempty
//     }

//     // Non può essere nel futuro
//     if date.After(time.Now().UTC()) {
//         return false
//     }

//     // Età minima 13 anni (esempio)
//     minAge := time.Now().AddDate(-13, 0, 0).UTC()
//     return date.Before(minAge)
// })
