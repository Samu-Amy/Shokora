package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Samu-Amy/Shokora/internal/api"
	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/store"
	"go.uber.org/zap"
)

func newTestApp(t *testing.T, showLogs bool) *api.App {
	t.Helper()

	mockStore := store.NewMockStore()
	var logger *zap.SugaredLogger

	testAuthenticator := &auth.TestAuthenticator{}

	if showLogs {
		logger = zap.Must(zap.NewProduction()).Sugar()
	} else {
		logger = zap.NewNop().Sugar()
	}

	return api.NewMockApp(&mockStore, logger, testAuthenticator)
}

func execureRequest(req *http.Request, router http.Handler) *httptest.ResponseRecorder {
	reqRec := httptest.NewRecorder() // Request Recorder
	router.ServeHTTP(reqRec, req)
	return reqRec
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("expected the response code to be %d but we got %d", expected, actual)
	}
}
