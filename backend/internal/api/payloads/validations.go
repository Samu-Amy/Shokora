package payloads

// import "github.com/go-playground/validator/v10"

// var Validate *validator.Validate //! andrebbe usato come payloads.Validate (oppure salvato come "Validator" dentro app ed usato come app.Validator)

// func InitValidators() {
// 	Validate = validator.New(validator.WithRequiredStructEnabled())

// 	Validate.RegisterValidation("password", passwordValidator)
// }

// - Auth -

// Password
var commonPasswords = map[string]struct{}{
	"password": {},
	"12345678": {},
	"qwerty":   {},
}

func IsCommonPassword(password string) bool {
	_, found := commonPasswords[password]
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
