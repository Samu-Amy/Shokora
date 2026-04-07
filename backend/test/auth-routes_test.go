package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	authservice "github.com/Samu-Amy/Shokora/internal/service/auth"
)

// TODO: verifica i vari casi e l'errore ottenuto (per ogni caso) - per controllare anche il parsing degli errori

func TestRegisterUserRoute(t *testing.T) {
	t.Run("should register all users", func(t *testing.T) {
		logRes := false

		for i := range min(routesTestsNum, len(validEmails)) { // Emails should be unique

			// Payload
			firstName := randomFrom(validFirstNames) // TODO: fare casi "fissi" hardcoded, invece che casuali (quelli già testati da validation e fuzz tests)
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

		clearTestDB(db)

		query := `
			INSERT INTO users (google_id, first_name, last_name, email, password, birthday, is_verified)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`

		// Create users in db
		for i := range min(routesTestsNum, len(validEmails)) {

			firstName := randomFrom(validFirstNames)
			lastName := randomFrom(validLastNames)
			strBirthday := randomFrom(validBirthdays)
			email := validEmails[i]
			pssw := randomFrom(validPasswords)

			hashedPssw, err := testService.Auth.HashPassword(strings.TrimSpace(pssw))
			if err != nil {
				t.Errorf("Error hashing password: %v", err)
			}

			var birthday time.Time
			if strBirthday != "" {
				birthday, err = authservice.ConvertBirthdayToTime(strings.TrimSpace(strBirthday))
				if err != nil {
					t.Errorf("Error converting birthday: %v", err)
				}
			}

			_, err = db.ExecContext(
				context.Background(),
				query,
				nil,
				firstName,
				lastName,
				email,
				hashedPssw,
				birthday,
				false,
			)

			if err != nil {
				t.Errorf("Error db: %v", err)
			}
		}

		// Test register handler
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

	// TODO: crea manualmente gli user prima (non è detto che siano presenti quelli usati da register) -> magari fai package per seeding

	// TODO: test login user esistenti e user non esistenti (usando tutti i dati mescolati, in modo da usare email di utenti con password di altri)
	// })
}
