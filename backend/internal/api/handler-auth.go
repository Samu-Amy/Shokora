package api

import (
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
	"github.com/go-chi/chi/v5"
)

// ----- REGISTER -----

// TODO: l'accesso con google (fai handler apposta) sostituisce solo la parte di autenticazione (login e register (in questo caso fornisce già la verifica della mail, settata a true)), poi la gestione di accesso e sessione è gestita dal mio sistema (?)

/*
Registers the user:
  - creates user account in db
  - handles email verification (creating tokens and sending email)
  - handles auth (creating tokens and setting cookies)
*/
func (app *App) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload data
	var payload payloads.RegisterUserReq

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Register user
	registerUserRes, authTokensDto, err := app.service.Auth.RegisterUser(ctx, payload)
	if err != nil {
		app.parseError(w, r, err)
		return
	}

	// Check cookies data
	if authTokensDto == nil {
		app.internalServerError(w, r, domerrors.ErrInternalError)
		return
	}

	// Set cookies
	app.setAuthCookies(w, *authTokensDto)

	// TODO: ricordati di scrivere di controllare nello spam (aggiungere timer al tasto per reinviare la mail (?))
	// TODO: se mail non inviata, dire di riprovare più tardi? -> l'utente può accedere ma non può ordinare (ha come opzioni di re-inviare la mail di verifica oppure eliminare l'account (e il token))

	//* Return user and verificationID
	if err := app.jsonResponse(w, http.StatusCreated, registerUserRes); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// ----- LOGIN/LOGOUT -----

func (app *App) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload data
	var payload payloads.LoginUserReq

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Log in user
	loginUserRes, authTokensDto, err := app.service.Auth.LoginUser(ctx, payload)
	if err != nil {
		app.parseError(w, r, err)
		return
	}

	// Set cookies if no 2fa
	if loginUserRes.User != nil && authTokensDto != nil {
		app.setAuthCookies(w, *authTokensDto)
	}

	//* Return user
	if err := app.jsonResponse(w, http.StatusCreated, loginUserRes); err != nil { // TODO: lato frontend bisognerà gestire i casi (es. call route per verifica)
		app.internalServerError(w, r, err)
		return
	}
}

// TODO: verifica che funzioni correttamente
func (app *App) logoutUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get sessionId from context (auth middleware)
	sessionId, ok := r.Context().Value(sessionIdCtx).(int64)
	if !ok {
		app.unauthorizedError(w, r, domerrors.ErrUnauthorized)
		return
	}

	// Logout (delete session)
	err := app.service.Auth.LogoutUser(ctx, sessionId)
	if err != nil {
		app.parseError(w, r, err)
		return
	}

	// Delete cookies
	app.clearAuthCookies(w)

	w.WriteHeader(http.StatusNoContent)
}

// ----- EMAIL VERIFICATION -----

func (app *App) verifyEmailWithMagicLinkHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get "token" param
	plainToken := chi.URLParam(r, verificationTokenParam)

	// Verify
	if err := app.service.Auth.VerifyEmailWithMagicLink(ctx, plainToken); err != nil {
		app.parseError(w, r, err) // TODO: nel FRONTEND dire che "non è valido o è scaduto" (non specificare quale dei due)
		return
	}

	//* No content
	w.WriteHeader(http.StatusNoContent)
}

func (app *App) verifyEmailWithOtpHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload data
	var payload payloads.OTPVerificationReq

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Verify
	if err := app.service.Auth.VerifyEmailWithOTP(ctx, &payload); err != nil {
		app.parseError(w, r, err)
		return
	}

	//* No content
	w.WriteHeader(http.StatusNoContent)
}

// ----- PASSWORD RESET -----

// - Request -

func (app *App) requestPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload data
	var payload payloads.RequestPasswordResetReq

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Request pasword reset
	if err := app.service.Auth.RequestPasswordReset(ctx, payload.Email); err != nil {
		app.parseError(w, r, err)
		return
	}

	//* No content
	w.WriteHeader(http.StatusNoContent)
}

// - Verification -

func (app *App) verifyPasswordResetWithMagicLinkHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: versione logged (usa user Id e magari richiede la password "vecchia") e versione non logged (la quale richiede l'email per poter verificare l'otp (in questo caso legato a email invece che user Id))

	ctx := r.Context()

	// Get "token" param
	plainToken := chi.URLParam(r, verificationTokenParam)

	// Verify
	plainResetSessionToken, err := app.service.Auth.VerifyPasswordResetWithMagicLink(ctx, plainToken)
	if err != nil {
		app.parseError(w, r, err)
		return
	}

	//* Return reset session token
	if err := app.jsonResponse(w, http.StatusCreated, plainResetSessionToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *App) verifyPasswordResetWithOtpHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: versione logged (usa user Id e magari richiede la password "vecchia") e versione non logged (la quale richiede l'email per poter verificare l'otp (in questo caso legato a email invece che user Id))

	ctx := r.Context()

	// Get payload data
	var payload payloads.OTPVerificationReq

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Verify
	plainResetSessionToken, err := app.service.Auth.VerifyPasswordResetWithOTP(ctx, &payload)
	if err != nil {
		app.parseError(w, r, err)
		return
	}

	//* Return reset session token
	if err := app.jsonResponse(w, http.StatusCreated, plainResetSessionToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// - Reset -

func (app *App) resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload data
	var payload payloads.ResetPasswordReq

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := app.service.Auth.ResetPassword(ctx, &payload)

	// TODO: loggare l'utente se non loggato (?)
}

// ----- TWO FACTOR AUTH -----

func (app *App) verifyTwoFactorAuthWithOtpHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload data
	var payload payloads.OTPVerificationReq

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Verify
	authTokensDto, err := app.service.Auth.TwoFactorAuthWithOTP(ctx, &payload)
	if err != nil {
		app.parseError(w, r, err)
		return
	}

	// Check cookies data
	if authTokensDto == nil {
		app.internalServerError(w, r, domerrors.ErrInternalError)
		return
	}

	// Set cookies
	app.setAuthCookies(w, *authTokensDto)

	// //* No content
	w.WriteHeader(http.StatusNoContent)
}
