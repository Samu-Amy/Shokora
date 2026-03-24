package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
)

func TestRegisterUser(t *testing.T) {
	t.Run("should register all users", func(t *testing.T) {
		logRes := false

		// Payload
		registerUserRequest := payloads.RegisterUserReq{
			UserDataReq: payloads.UserDataReq{
				FirstName: "John",
				LastName:  "Snow",
			},
			EmailFieldReq: payloads.EmailFieldReq{
				Email: "john.snow@gmail.com",
			},
			DoublePasswordFieldReq: payloads.DoublePasswordFieldReq{
				Password:             "RobStark",
				PasswordConfirmation: "RobStark",
			},
		}

		// Request
		w := makeRequestWithPayload(t, testRouter, "POST", "/api/v1/auth/user", registerUserRequest)

		// Check the result
		checkResponseCode(t, w, http.StatusBadRequest)

		// Log response body
		if logRes {
			logResBody(t, w)
		}
	})

	t.Run("should give badRequest errors because of invalid data in payload", func(t *testing.T) {
		logRes := false

		// TODO: sistema validazione request payloads

		// TODO: fai payloads con field non validi

		// Payload
		registerUserRequest := payloads.RegisterUserReq{
			UserDataReq: payloads.UserDataReq{
				FirstName: "John&", // Symbols in first/last name
				LastName:  "%Snow",
				Birthday:  time.Now(), // Date not valid
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
		checkErrorMessage(t, w, "bad_request")

		// Log response body
		if logRes {
			logResBody(t, w)
		}
	})

	t.Run("should give badRequest errors because of common password", func(t *testing.T) {
		logRes := false

		// Payload
		registerUserRequest := payloads.RegisterUserReq{
			UserDataReq: payloads.UserDataReq{
				FirstName: "John&",
				LastName:  "Snow",
			},
			EmailFieldReq: payloads.EmailFieldReq{
				Email: "john.snow@gmail.doc",
			},
			DoublePasswordFieldReq: payloads.DoublePasswordFieldReq{
				Password:             "admin12345678",
				PasswordConfirmation: "admin12345678",
			},
		}

		// Request
		w := makeRequestWithPayload(t, testRouter, "POST", "/api/v1/auth/user", registerUserRequest)

		// Check the result
		checkResponseCode(t, w, http.StatusBadRequest)
		checkErrorMessage(t, w, "common_password")

		// Log response body
		if logRes {
			logResBody(t, w)
		}
	})

	t.Run("should give badRequest errors because of different password", func(t *testing.T) {
		logRes := false

		// TODO: fai controllo (che lo user abbia messo la password che vuole mettendola giusta due volte)?

		// Payload
		registerUserRequest := payloads.RegisterUserReq{
			UserDataReq: payloads.UserDataReq{
				FirstName: "John&",
				LastName:  "Snow",
			},
			EmailFieldReq: payloads.EmailFieldReq{
				Email: "john.snow@gmail.doc",
			},
			DoublePasswordFieldReq: payloads.DoublePasswordFieldReq{
				Password:             "RobStark",
				PasswordConfirmation: "StarkRob",
			},
		}

		// Request
		w := makeRequestWithPayload(t, testRouter, "POST", "/api/v1/auth/user", registerUserRequest)

		// Check the result
		checkResponseCode(t, w, http.StatusBadRequest)
		checkErrorMessage(t, w, "different_password")

		// Log response body
		if logRes {
			logResBody(t, w)
		}
	})
}

func TestLoginUser(t *testing.T) {
	// t.Run("should not allow unauthenticated requests", func (t *testing.T) {
	// TODO: test login user esistenti e user non esistenti (usando tutti i dati mescolati, in modo da usare email di utenti con password di altri)
	// })
}
