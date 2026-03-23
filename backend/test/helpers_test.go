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

func checkResponseCode(t *testing.T, expected int, w *httptest.ResponseRecorder) {
	t.Helper()

	if expected != w.Code {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

func logResBody(t *testing.T, w *httptest.ResponseRecorder) {
	if w.Body.Len() > 0 {
		t.Logf("response body: %s", w.Body.String())
	}
}
