package api

import (
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/store"
)

// - REGISTER -

// TODO: aggiungi data di nascita (opzionale)
type RegisterUserPayload struct {
	FirstName string `json:"first_name" validate:"required,max=125"`
	LastName  string `json:"last_name" validate:"required,max=125"`
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=6,max=72"`
	// BirthDate string `json:"birth_date,omitempty" validate:"omitempty"` // TODO: togliere omitempty (?)
}

func (app *App) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload data
	var payload RegisterUserPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Create product from payload data
	user := &store.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
	}

	// Hash password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Generate Token

	// Create product
	if err := app.store.User.CreateAndSendVerification(ctx, user, "token-123", app.config.Mail.EmailVerificationTokenExp); err != nil {
		app.parseError(w, r, err)
		return
	}

	//* Return user
	if err := app.jsonResponse(w, http.StatusCreated, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
