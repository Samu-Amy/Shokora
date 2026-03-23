package main

import (
	"net/http"
	"testing"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
)

func TestRegisterUser(t *testing.T) {
	t.Run("should register all users", func(t *testing.T) {
		// testService.Auth.RegisterUser()
	})

	t.Run("should give badRequest errors because of wrong payload", func(t *testing.T) {
		registerUserRequest := payloads.RegisterUserReq{
			UserDataReq: payloads.UserDataReq{
				FirstName: "John",
				LastName:  "Snow",
			},
			EmailFieldReq: payloads.EmailFieldReq{
				Email: "john.snow@gmail.com",
			},
			PasswordFieldReq: payloads.PasswordFieldReq{
				Password: "RobStark",
			},
		}

		w := makeRequestWithPayload(t, testRouter, "POST", "/api/v1/auth/user", registerUserRequest)

		// Check the result
		checkResponseCode(t, http.StatusBadRequest, w)

		// Log response body
		// t.Logf("Response body: %v", w.Body)
	})
}

func TestLoginUser(t *testing.T) {
	// t.Run("should not allow unauthenticated requests", func (t *testing.T) {
	// TODO: test login user esistenti e user non esistenti (usando tutti i dati mescolati, in modo da usare email di utenti con password di altri)
	// })
}
