package payloads

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func GetValidationErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "Questo campo è obbligatorio"

	case "min":
		return "Deve contenere contenere almeno " + fieldError.Param() + " caratteri"

	case "max":
		return "Deve contenere contenere al massimo " + fieldError.Param() + " caratteri"

	case "email":
		return "Email non valida"

	case "eqfield":
		return "Deve coincidere con " + strings.ToLower(fieldError.Param())

	case "nefield":
		return "Deve essere diversa da " + strings.ToLower(fieldError.Param())

	case "valid-name":
		return "Contiene caratteri non validi"

	case "no-edge-spaces":
		return "Non sono ammessi spazi e inizio e fine"

	case "valid-password":
		return "Contiene caratteri non validi"

	case "no-common-password":
		return "Password troppo comune, scegline una più sicura"

	default:
		return "Valore non valido"
	}
}

// Es. for language
// var messages = map[string]map[string]string{
//     "it": {
//         "required": "Questo campo è obbligatorio",
//         "min":      "Deve contenere almeno %s caratteri",
//         "max":      "Deve contenere al massimo %s caratteri",
//         ...
//     },
//     "en": {
//         "required": "This field is required",
//         "min":      "Must be at least %s characters",
//         ...
//     },
// }
