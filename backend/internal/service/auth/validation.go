package authservice

import "regexp"

// import "time"

// const minAge = 16

// func isAgeValid(birthDate time.Time) bool {
// 	// Basic check
// 	if birthDate.IsZero() || birthDate.UTC().After(time.Now().UTC()) {
// 		return false
// 	}

// 	// Age check: calculate when the user turn minAge
// 	birthday := time.Date(
// 		birthDate.Year()+minAge,
// 		birthDate.Month(),
// 		birthDate.Day(),
// 		0, 0, 0, 0, time.UTC,
// 	)

// 	return !time.Now().UTC().Before(birthday)
// }

func isNameValid(password string) bool {
	// No symbols
}

func isStrinValid(text string) bool {
	// TODO: controlla emojii
}

func isPasswordValid(password string) bool {
	var passwordRegexp = regexp.MustCompile(`^[\x21-\x7E]+$`) // TODO: attenzione non sono accettati caratteri tipo 'ò'
}
