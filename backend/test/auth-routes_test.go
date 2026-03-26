package main

import (
	"net/http"
	"testing"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
)

// TODO: verifica i vari casi e l'errore ottenuto (per ogni caso) - per controllare anche il parsing degli errori

func TestRegisterUserRoute(t *testing.T) {
	t.Run("should register all users", func(t *testing.T) {
		logRes := false

		// TODO: sistema mock data (minor numero di dati ma che copre il maggior numero di casi possibili, tanto per la validazione dei dati ci sono test apposta)

		for range 50 {
			// Payload
			registerUserReq := makeRegisterUserReq(
				randomFrom(validFirstNames),
				randomFrom(validLastNames),
				randomFrom(validBirthdays),
				randomFrom(validEmails),
				randomFrom(validPasswords),
				randomFrom(validPasswords),
			)

			// Request
			w := makeRequestWithPayload(t, testRouter, "POST", "/api/v1/auth/user", registerUserReq)

			// Check the result
			checkResponseCode(t, w, http.StatusNoContent)

			// Log response body
			if logRes {
				logResBody(t, w)
			}
		}
	})

	// Payload validation tested separately (here is just to test that works in the route)
	t.Run("should give badRequest errors because of invalid data in payload", func(t *testing.T) {
		logRes := false

		// TODO: fai payloads con field non validi

		// Payload
		registerUserRequest := payloads.RegisterUserReq{
			UserDataReq: payloads.UserDataReq{
				FirstName: "John&", // Symbols in first/last name
				LastName:  "%Snow",
			},
			EmailFieldReq: payloads.EmailFieldReq{
				Email: "john.snow@gmail.doc", // Invalid email
			},
			DoublePasswordFieldReq: payloads.DoublePasswordFieldReq{
				Password:             "Rob", // Password too short
				PasswordConfirmation: "Rob",
			},
		}

		// Request
		w := makeRequestWithPayload(t, testRouter, "POST", "/api/v1/auth/user", registerUserRequest)

		// Check the result
		checkResponseCode(t, w, http.StatusBadRequest)
		// checkErrorMessage(t, w, "bad_request") // TODO: dipende dai dati non validi (fai check generico oppure specifico ma con casi specifici (es. sapendo che in un caso la password non passa una certa validazione, aspettati quel'errore specifico))

		// Log response body
		if logRes {
			logResBody(t, w)
		}
	})
}

func TestLoginUserRoute(t *testing.T) {
	// t.Run("should not allow unauthenticated requests", func (t *testing.T) {
	// TODO: test login user esistenti e user non esistenti (usando tutti i dati mescolati, in modo da usare email di utenti con password di altri)
	// })
}
