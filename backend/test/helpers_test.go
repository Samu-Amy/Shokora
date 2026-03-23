package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeRequestWithPayload(t *testing.T, router http.Handler, method, route string, payload any) *httptest.ResponseRecorder {
	t.Helper() // indica che è un helper, utile per errori

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
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
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

func logResBody(t *testing.T, w *httptest.ResponseRecorder) {
	if w.Body.Len() > 0 {
		t.Logf("response body: %s", w.Body.String())
	}
}
