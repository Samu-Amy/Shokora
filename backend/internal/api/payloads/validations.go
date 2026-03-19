package payloads

import "strings"

// import "github.com/go-playground/validator/v10"

// var Validate *validator.Validate //! andrebbe usato come payloads.Validate (oppure salvato come "Validator" dentro app ed usato come app.Validator)

// func InitValidators() {
// 	Validate = validator.New(validator.WithRequiredStructEnabled())

// 	Validate.RegisterValidation("password", passwordValidator)
// }

// - Auth -

// Password
var commonPasswords = map[string]struct{}{
	"password":         {},
	"12345678":         {},
	"qwerty":           {},
	"admin":            {},
	"passwordpassword": {},
	"123456789012":     {},
	"qwertyuiop123":    {},
	"iloveyou123456":   {},
	"letmein1234567":   {},
	"welcome123456":    {},
	"admin12345678":    {},
	"monkey123456789":  {},
	"dragon12345678":   {},
	"shadow12345678":   {},
	"master12345678":   {},
	"football123456":   {},
	"princess123456":   {},
	"sunshine123456":   {},
	"baseball1234567":  {},
	"098765432109":     {},
	"123412341234":     {},
	"202312345678":     {},
	"202412345678":     {},
	"qwerty12345678":   {},
	"111111111111":     {},
	"999999999999":     {},
	"000000000000":     {},
	"654321654321":     {},
	"passwordpass123":  {},
	"passw0rdpass123":  {},
	"P@ssword123456":   {},
}

func IsCommonPassword(password string) bool {
	_, found := commonPasswords[strings.ToLower(password)]
	return found
}

// func passwordValidator(fl validator.FieldLevel) bool {
// 	pwd := fl.Field().String()

// 	hasLength := len(pwd) >= 12 && len(pwd) <= 72

// 	_, isCommon := commonPasswords[pwd]

// 	return hasLength && !isCommon
// }

// Birthdate

// TODO: fare validazione custom (tipo quella sotto)?

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
