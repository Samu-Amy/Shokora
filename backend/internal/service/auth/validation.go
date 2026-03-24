package authservice

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
