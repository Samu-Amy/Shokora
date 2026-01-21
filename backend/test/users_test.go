package main

import (
	"net/http"
	"testing"
)

// TODO: fai test più "ampi" con db di test ed elimina mock (da store, auth e funzioni aggiuntive in server)

// TODO: aggiungere spies (?)

func TestGetUserHandler(t *testing.T) {
	app := newTestApp(t, true)
	router := app.GetRouter()

	testToken, err := app.GetAuthenticator().GenerateToken(nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/users/42", nil)
		if err != nil {
			t.Fatal(err)
		}

		reqRec := execureRequest(req, router)

		checkResponseCode(t, http.StatusUnauthorized, reqRec.Code)
	})

	// TODO: da implementare (deve funzionare solo se userID == id dell'utente nell'auth)
	t.Run("should allow authenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/users/42", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		reqRec := execureRequest(req, router)

		checkResponseCode(t, http.StatusOK, reqRec.Code)
	})
}
