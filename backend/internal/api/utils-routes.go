package api

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/Samu-Amy/Shokora/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// - Context -

// Keys
type contextKey int

const (
	userCtx contextKey = iota
)

// Functions
func getUserFromContext(r *http.Request) (*store.User, bool) {
	user, ok := r.Context().Value(userCtx).(*store.User)
	return user, ok
}

// - Params -

func (app *App) getIdFromParam(r *http.Request, idParamName string) (int64, error) {
	idParam := chi.URLParam(r, idParamName)

	resourceId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return -1, err
	}

	return resourceId, nil
}

// - Token Generation -

func (app *App) generateHashedToken() (string, string) {
	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	return hashedToken, plainToken
}
