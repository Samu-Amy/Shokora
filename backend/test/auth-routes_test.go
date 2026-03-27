package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"testing"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
)

// TODO: verifica i vari casi e l'errore ottenuto (per ogni caso) - per controllare anche il parsing degli errori

func TestRegisterUserRoute(t *testing.T) {
	t.Run("should register all users", func(t *testing.T) {
		logRes := false

		// TODO: sistema mock data (minor numero di dati ma che copre il maggior numero di casi possibili, tanto per la validazione dei dati ci sono test apposta)

		for i := range min(routesTestsNum, len(validEmails)) { // Emails should be unique

			// Payload
			firstName := randomFrom(validFirstNames)
			lastName := randomFrom(validLastNames)
			birthday := randomFrom(validBirthdays)
			email := validEmails[i]
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

			// Check important data
			if res.User.IsVerified || res.User.Role != 0 {
				t.Fatal("user created as verified or with role != 0")
			}

			// if res.User.FirstName != strings.TrimSpace(firstName) || res.User.LastName != strings.TrimSpace(lastName) || res.User.Email != strings.TrimSpace(email) {
			// 	t.Fatal("wrong user data in payload")
			// }
		}
	})

	t.Run("should give badRequest error because of email already in db", func(t *testing.T) {
		logRes := true

		for i := range min(routesTestsNum, len(validEmails)) {

			// Payload
			firstName := randomFrom(validFirstNames)
			lastName := randomFrom(validLastNames)
			birthday := randomFrom(validBirthdays)
			email := validEmails[i]
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
			checkResponseCode(t, w, http.StatusBadRequest)
			checkErrorMessage(t, w, "duplicate_email")

			// Log response body
			if logRes {
				logResBody(t, w)
			}
		}
	})

	// Payload validation tested separately (here is just to test that works in the route)
	t.Run("should give badRequest errors because of invalid data in payload", func(t *testing.T) {
		logRes := true

		for i := range min(routesTestsNum, len(notValidEmails)) { // Emails should be unique

			// Payload
			firstName := randomFrom(notValidFirstNames)
			lastName := randomFrom(notValidLastNames)
			birthday := randomFrom(notValidBirthdays)
			email := notValidEmails[i]
			pssw := randomFrom(notValidPasswords)
			pssw2 := pssw

			if rand.Float32() < 0.5 {
				pssw2 = notValidPasswords[rand.Intn(len(notValidPasswords))] // Pick a random password (not using the customRand with the same seed as others)
			}

			registerUserReq := makeRegisterUserReq(
				firstName,
				lastName,
				birthday,
				email,
				pssw,
				pssw2,
			)

			// Request
			w := makeRequestWithPayload(t, testRouter, "POST", "/api/v1/auth/user", registerUserReq)

			// Check the result
			checkResponseCode(t, w, http.StatusBadRequest)

			// Log response body
			if logRes {
				logResBody(t, w)
			}
		}
	})
}

func TestLoginUserRoute(t *testing.T) {
	// t.Run("should not allow unauthenticated requests", func (t *testing.T) {
	// TODO: test login user esistenti e user non esistenti (usando tutti i dati mescolati, in modo da usare email di utenti con password di altri)
	// })
}
