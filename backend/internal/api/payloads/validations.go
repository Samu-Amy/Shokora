package payloads

// validate := validator.New(validator.WithRequiredStructEnabled())
var validate = validator.New()

// - Auth -

// Password
var commonPasswords = map[string]struct{}{
    "password": {},
    "12345678": {},
    "qwerty":   {},
}

validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
    pwd := fl.Field().String()

	hasLength := len(pwd) >= 12 && len(pwd) <= 72
	
	_, isCommon := commonPasswords[pwd]
	
    return hasLength && !isCommon
})

// hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pwd)
// hasLower := regexp.MustCompile(`[a-z]`).MatchString(pwd)
// hasNumber := regexp.MustCompile(`[0-9]`).MatchString(pwd)
// hasSpecial := regexp.MustCompile(`[!@#~$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(pwd)

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
