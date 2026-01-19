package api

import (
	"fmt"
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/mailer"
	"github.com/Samu-Amy/Shokora/internal/store"
	"github.com/go-chi/chi/v5"
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

	// Hash and set password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Generate Token
	hashedToken, plainToken := app.generateHashedToken()

	// Create User
	if err := app.store.User.CreateUserAndSendVerification(ctx, user, hashedToken, app.config.Mail.EmailVerificationTokenExp); err != nil {
		app.parseError(w, r, err)
		return
	} // TODO: gestire meglio (verificare scadenza token, se scaduto cosa si fa?)

	activationURL := fmt.Sprintf("%s/verify-email/%s", app.config.FrontEndURL, plainToken)

	vars := struct {
		Name          string
		ActivationURL string
	}{
		Name:          user.FirstName,
		ActivationURL: activationURL,
	}

	isProdEnv := app.config.Env == "prod"

	// Send email
	err := app.mailer.SendEmail(ctx, mailer.EmailVerificationTemplate, user.FirstName, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending welcome email", "error", err)

		// Rollback user creation
		if err := app.store.User.DeleteUserAndEmailVerificationToken(ctx, user.Id); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}

		app.internalServerError(w, r, err) // TODO: evitare di eliminare user e token e dire di riprovare più tardi -> l'utente può accedere ma non può ordinare (ha come opzioni di re-inviare la mail di verifica oppure eliminare l'account (e il token))
		return
	}

	app.logger.Infow("Email sent successfully")

	// TODO: ricordati di scrivere di controllare nello spam

	//* Return user
	if err := app.jsonResponse(w, http.StatusCreated, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// - EMAIL VERIFICATION -

func (app *App) verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := chi.URLParam(r, "token")

	if err := app.store.User.VerifyEmail(ctx, token); err != nil {
		app.parseError(w, r, err)
		return
	}

	//* No content
	w.WriteHeader(http.StatusNoContent)
}
