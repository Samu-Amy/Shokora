package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
)

// TODO: verifica i vari casi e l'errore ottenuto (per ogni caso) - per controllare anche il parsing degli errori

func TestRegisterUserRoute(t *testing.T) {
	t.Run("should register all users", func(t *testing.T) {
		logRes := true

		// TODO: sistema mock data (minor numero di dati ma che copre il maggior numero di casi possibili, tanto per la validazione dei dati ci sono test apposta)

		for range routesTestsNum {

			// Payload
			firstName := randomFrom(validFirstNames)
			lastName := randomFrom(validLastNames)
			birthday := randomFrom(validBirthdays)
			email := randomFrom(validEmails)
			pssw := randomFrom(validPasswords)

			registerUserReq := makeRegisterUserReq(
				firstName,
				lastName,
				birthday,
				email,
				pssw,
				pssw,
			)

			// Request
			w := makeRequestWithPayload(t, testRouter, "POST", "/api/v1/auth/user", registerUserReq)

			// Check the result
			checkResponseCode(t, w, http.StatusCreated)

			// Log response body
			if logRes {
				logResBody(t, w)
			}

			// - Check response -

			var res payloads.RegisterUserRes

			err := json.Unmarshal(w.Body.Bytes(), &res)
			if err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}

			if res.User.FirstName != strings.TrimSpace(firstName) || res.User.LastName != strings.TrimSpace(lastName) || res.User.Email != strings.TrimSpace(email) {
				t.Fatal("wrong user data in payload")
			}
		}
	})

	// Payload validation tested separately (here is just to test that works in the route)
	// t.Run("should give badRequest errors because of invalid data in payload", func(t *testing.T) {
	// 	logRes := false

	// 	// TODO: fai payloads con field non validi

	// 	// Payload
	// 	registerUserRequest := payloads.RegisterUserReq{
	// 		UserDataReq: payloads.UserDataReq{
	// 			FirstName: "John&", // Symbols in first/last name
	// 			LastName:  "%Snow",
	// 		},
	// 		EmailFieldReq: payloads.EmailFieldReq{
	// 			Email: "john.snow@gmaildoc", // Invalid email
	// 		},
	// 		DoublePasswordFieldReq: payloads.DoublePasswordFieldReq{
	// 			Password:             "Rob",         // Password too short
	// 			PasswordConfirmation: "Rob1231sad3", // Different password
	// 		},
	// 	}

	// 	// Request
	// 	w := makeRequestWithPayload(t, testRouter, "POST", "/api/v1/auth/user", registerUserRequest)

	// 	// Check the result
	// 	checkResponseCode(t, w, http.StatusBadRequest)
	// 	// checkErrorMessage(t, w, "bad_request") // TODO: dipende dai dati non validi (fai check generico oppure specifico ma con casi specifici (es. sapendo che in un caso la password non passa una certa validazione, aspettati quel'errore specifico))

	// 	// Log response body
	// 	if logRes {
	// 		logResBody(t, w)
	// 	}
	// })
}

func TestLoginUserRoute(t *testing.T) {
	// t.Run("should not allow unauthenticated requests", func (t *testing.T) {
	// TODO: test login user esistenti e user non esistenti (usando tutti i dati mescolati, in modo da usare email di utenti con password di altri)
	// })
}
