package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Samu-Amy/Shokora/internal/api/payload"
	"github.com/Samu-Amy/Shokora/internal/mailer"
	"github.com/Samu-Amy/Shokora/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ----- REGISTER -----

func (app *App) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload data
	var payload payload.RegisterUserPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Create user from payload data
	user := &store.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
	}

	// Hash and set password
	if err := user.Password.Set(payload.Password); err != nil { // TODO: sistema
		app.internalServerError(w, r, err)
		return
	}

	// Generate Token
	hashedToken, plainToken := app.generateHashedToken() // TODO: sistema

	// Create User
	if err := app.store.User.CreateUserAndSendVerification(ctx, user, hashedToken, app.config.Auth.MagicLink.Exp); err != nil {
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

// ----- EMAIL VERIFICATION -----

func (app *App) verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := chi.URLParam(r, "token")

	if err := app.store.User.VerifyEmail(ctx, token); err != nil {
		app.parseError(w, r, err)
		return
	}

	// TODO: aggiungi creazione ed invio token (togliere createTokenHandler?) - crea anche refresh token (?)

	//* No content
	w.WriteHeader(http.StatusNoContent) // TODO: setta http-only cookies con token e invia user
}

// ----- TOKENS -----

func (app *App) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload data
	var payload payload.CreateUserTokenPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Fetch the user (check if the user exist)
	user, err := app.store.User.GetByEmail(ctx, payload.Email)
	if err != nil {
		app.parseError(w, r, err) // TODO: non dire se l'email esiste o meno
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword(user.Password.Hash, []byte(payload.Password))
	if err != nil {
		app.unauthorizedError(w, r, err)
		return
	}

	// Generate tokens (and add claims)
	claims := jwt.MapClaims{
		"sub": user.Id, // subject
		"exp": time.Now().Add(app.config.Auth.Token.AccessTokenExp).Unix(),
		"iat": time.Now().Unix(),              // issued at
		"nbf": time.Now().Unix(),              // not before time
		"iss": app.config.Auth.Token.Issuer,   // issuer
		"aud": app.config.Auth.Token.Audience, // audience
	}

	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// TODO: setta cookie invece che inviarlo come payload (?)

	//* Send token to the client
	if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *App) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: opzioni per cookie (da verificare)
	// HttpOnly: true
	// Secure: true
	// SameSite: Strict (refresh) / Lax (access)
	// Path: /auth/refresh (per refresh token)
}
