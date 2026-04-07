package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"testing"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/auth"
	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
)

// TODO: verifica i vari casi e l'errore ottenuto (per ogni caso) - per controllare anche il parsing degli errori

func TestRegisterUserRoute(t *testing.T) {
	t.Run("should register all users", func(t *testing.T) {
		logRes := false

		clearTestDB(db)

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

			var res APIResponse[payloads.RegisterUserRes]

			err := json.Unmarshal(w.Body.Bytes(), &res)
			if err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}

			// Check important data
			if res.Data.User.IsVerified || res.Data.User.Role != 0 {
				t.Fatal("user created as verified or with role != 0")
			}

			if res.Data.User.FirstName != strings.TrimSpace(firstName) || res.Data.User.LastName != strings.TrimSpace(lastName) || res.Data.User.Email != strings.TrimSpace(email) {
				t.Fatal("wrong user data in payload")
			}
		}
	})

	t.Run("should give badRequest error because of email already in db", func(t *testing.T) {
		logRes := false

		clearTestDB(db)
		seedUsers(t, db)

		// Test register handler
		for i := range min(seedUserNum, routesTestsNum, len(validEmails)) {

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
			checkErrorMessage(t, w, domerrors.ErrDuplicateEmail.Error())

			// Log response body
			if logRes {
				logResBody(t, w)
			}
		}
	})

	// Payload validation tested separately (here is just to test that works in the route)
	t.Run("should give badRequest errors because of invalid data in payload", func(t *testing.T) {
		logRes := false

		for i := range min(routesTestsNum, len(notValidEmails)) { // Emails should be unique

			// Payload
			firstName := randomFrom(notValidFirstNames)
			lastName := randomFrom(notValidLastNames)
			birthday := randomFrom(notValidBirthdays)
			email := notValidEmails[i]
			pssw := randomFrom(notValidPasswords)
			pssw2 := pssw

			if rand.Float32() < 0.5 { // ~ 50% use PasswordConfirmation with different (wrong) value
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
	t.Run("should authenticate users", func(t *testing.T) {
		logRes := false

		clearTestDB(db)
		seedUsers(t, db)

		// Test login handler
		for i := range min(seedUserNum, routesTestsNum, len(validEmails)) {

			// Payload
			email := validEmails[i]
			pssw := validPasswords[i]

			loginUserReq := makeLoginUserReq(
				email,
				pssw,
			)

			// Request
			w := makeRequestWithPayload(t, testRouter, "GET", "/api/v1/auth/user", loginUserReq)

			// Check the result
			checkResponseCode(t, w, http.StatusOK)

			// Log response body
			if logRes {
				logResBody(t, w)
			}

			// - Check response -

			var res APIResponse[payloads.LoginUserRes]

			err := json.Unmarshal(w.Body.Bytes(), &res)
			if err != nil {
				t.Fatalf("failed to unmarshal response body: %v", err)
			}

			// Check important data
			if res.Data.VerificationId == nil {
				if res.Data.User == nil {
					t.Fatal("user not sent") // TODO: fix this error (response body: {"data":{"is_email_sent":false}})
				}

				if res.Data.User.IsVerified != userVerified[i] || res.Data.User.Role != userRoles[i] {
					t.Fatal("user has wrong is_verified or user_role")
				}

				if res.Data.User.FirstName != strings.TrimSpace(validFirstNames[i]) || res.Data.User.LastName != strings.TrimSpace(validLastNames[i]) {
					t.Fatal("user has wrong first or last name")
				}
			} else {

				verificationType := getVerificationType(t, res.Data.VerificationId)

				if verificationType == auth.EmailVerification && res.Data.User == nil {
					t.Fatal("user not sent while email verification requested")

				} else if verificationType == auth.TwoFactorAuth && res.Data.User != nil {
					t.Fatal("user sent while 2FA active")
				}
			}
		}
	})

	t.Run("should not authenticate users bacause of emails not exists in db", func(t *testing.T) {
		logRes := false

		clearTestDB(db)

		// Test login handler
		for i := range min(seedUserNum, routesTestsNum, len(validEmails)) {

			// Payload
			email := validEmails[i]
			pssw := validPasswords[i]

			loginUserReq := makeLoginUserReq(
				email,
				pssw,
			)

			// Request
			w := makeRequestWithPayload(t, testRouter, "GET", "/api/v1/auth/user", loginUserReq)

			// Check the result
			checkResponseCode(t, w, http.StatusBadRequest)
			checkErrorMessage(t, w, domerrors.ErrInvalid.Error())

			// Log response body
			if logRes {
				logResBody(t, w)
			}
		}
	})

	t.Run("should not authenticate users bacause of wrong passwords", func(t *testing.T) {
		logRes := false

		clearTestDB(db)
		seedUsers(t, db)

		// Test login handler
		for i := range min(seedUserNum, routesTestsNum, len(validEmails)) {

			// Payload
			email := validEmails[i]

			// Get a different index for the password
			randomIndex := customRand.Intn(len(validPasswords))
			for randomIndex == i {
				randomIndex = customRand.Intn(len(validPasswords))
			}

			pssw := validPasswords[randomIndex]

			loginUserReq := makeLoginUserReq(
				email,
				pssw,
			)

			// Request
			w := makeRequestWithPayload(t, testRouter, "GET", "/api/v1/auth/user", loginUserReq)

			// Check the result
			checkResponseCode(t, w, http.StatusBadRequest)
			checkErrorMessage(t, w, domerrors.ErrInvalid.Error())

			// Log response body
			if logRes {
				logResBody(t, w)
			}
		}
	})
}
