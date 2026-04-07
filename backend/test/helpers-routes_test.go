package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
)

// Res Data
type APIResponse[T any] struct {
	Data T `json:"data"`
}

// - Routes -

func makeRequestWithPayload(t *testing.T, router http.Handler, method, route string, payload any) *httptest.ResponseRecorder {
	t.Helper() // When calling t.Fatal, t.Error, ecc. it reports the line in the test instead of the one in this function

	// Marshal in JSON
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	// Crea richiesta HTTP
	req := httptest.NewRequest(method, route, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Recorder
	w := httptest.NewRecorder()

	// Chiama il router
	router.ServeHTTP(w, req)

	return w
}

func checkResponseCode(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if expected != w.Code {
		t.Fatalf("expected %d, got %d, body: %s", expected, w.Code, w.Body.String())
	}
}

func checkErrorMessage(t *testing.T, w *httptest.ResponseRecorder, expected string) {
	t.Helper()

	type errorPayload struct {
		Error string `json:"error"`
	}

	var res errorPayload

	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if res.Error != expected {
		t.Fatalf("expected error message %q, got %q", expected, res.Error)
	}
}

// func checkResBody[T any](t *testing.T, w *httptest.ResponseRecorder, expected T) {
// 	t.Helper()

// 	// TODO: fai (con generics?) - però ci sono dati che non so (tipo )

// 	var res T

// 	err := json.Unmarshal(w.Body.Bytes(), &res)
// 	if err != nil {
// 		t.Fatalf("failed to unmarshal response body: %v", err)
// 	}

// 	if res != expected {
// 		t.Fatalf("expected error message %q, got %q", expected, res.Error)
// 	}
// }

func logResBody(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()

	if w.Body.Len() > 0 {
		t.Logf("response body: %s", w.Body.String())
	}
}

// - Validation -
func parseValidationErr(t *testing.T, err error) validator.ValidationErrors {
	t.Helper()

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		t.Fatalf("unexpected error type: %v", err)
	}

	return validationErrors
}
