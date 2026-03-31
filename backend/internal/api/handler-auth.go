package api

import (
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/auth"
	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
	user_repo "github.com/Samu-Amy/Shokora/internal/store/user"
	"github.com/go-chi/chi/v5"
)

// TODO: l'accesso con google (fai handler apposta) sostituisce solo la parte di autenticazione (login e register (in questo caso fornisce già la verifica della mail, settata a true)), poi la gestione di accesso e sessione è gestita dal mio sistema (?)

//! TODO: nelle condizioni da accettare indica che l'utente dichiara di avere l'età minima per potersi registrare (e per poter fare acquisti)

// ----- REGISTER -----

func (app *App) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload data
	var payload payloads.RegisterUserReq // TODO: nel FRONTEND chiedi solo giorno e mese (compleanno) e non anno (non data di nascita) - evita anche emoji ed altro (nonostante i controlli qua)

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := app.dataValidator.Struct(payload); err != nil {
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

// ----- LOGIN -----

func (app *App) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload data
	var payload payloads.LoginUserReq

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := app.dataValidator.Struct(payload); err != nil {
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

func (app *App) googleHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	// TODO: controlla email esistente

	//* Return user
	// if err := app.jsonResponse(w, http.StatusCreated, loginUserRes); err != nil { // TODO: lato frontend bisognerà gestire i casi (es. call route per verifica)
	// 	app.internalServerError(w, r, err)
	// 	return
	// }
}

func (app *App) googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	//* Return user
	// if err := app.jsonResponse(w, http.StatusCreated, loginUserRes); err != nil { // TODO: lato frontend bisognerà gestire i casi (es. call route per verifica)
	// 	app.internalServerError(w, r, err)
	// 	return
	// }
}

// ----- LOGOUT -----

func (app *App) logoutUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get sessionId from context (auth middleware)
	sessionId, ok := ctx.Value(sessionIdCtx).(int64)
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

	//* No content
	w.WriteHeader(http.StatusNoContent)
}

// ----- GET USER -----

func (app *App) getCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from context (auth middleware)
	user, ok := ctx.Value(userCtx).(*user_repo.User)
	if !ok {
		app.unauthorizedError(w, r, domerrors.ErrUnauthorized)
		return
	}

	userRes := payloads.ToUserRes(*user)

	//* Return user
	if err := app.jsonResponse(w, http.StatusOK, userRes); err != nil {
		app.internalServerError(w, r, err)
		return
	}
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
	if err := app.dataValidator.Struct(payload); err != nil {
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

// - Verification -

func (app *App) verifyPasswordResetWithMagicLinkHandler(w http.ResponseWriter, r *http.Request) { // TODO: aggiungi controllo nel service per non far cambiare la password a chi ha l'accesso con google
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
	ctx := r.Context()

	// Get payload data
	var payload payloads.OTPVerificationReq

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := app.dataValidator.Struct(payload); err != nil {
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
	if err := app.dataValidator.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Reset password
	if err := app.service.Auth.ResetPassword(ctx, &payload); err != nil {
		app.parseError(w, r, err)
		return
	}

	// TODO: loggare l'utente (?)

	//* No content
	w.WriteHeader(http.StatusNoContent)
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
	if err := app.dataValidator.Struct(payload); err != nil {
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

	//* No content
	w.WriteHeader(http.StatusNoContent)
}

// ------ RESEND VERIFICATION -----

// - Email verification -

func (app *App) resendEmailVerificationHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// Get payload data
	var payload payloads.SendVerificationReq

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := app.dataValidator.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Request pasword reset
	if err := app.service.Auth.SendVerification(ctx, auth.EmailVerification, payload.Email); err != nil {
		app.parseError(w, r, err)
		return
	}

	//* No content
	w.WriteHeader(http.StatusNoContent)
}

// - Password Reset -

func (app *App) sendPasswordResetHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// Get payload data
	var payload payloads.SendVerificationReq

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := app.dataValidator.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Request pasword reset
	if err := app.service.Auth.SendVerification(ctx, auth.PasswordReset, payload.Email); err != nil {
		app.parseError(w, r, err)
		return
	}

	//* No content
	w.WriteHeader(http.StatusNoContent)
}

// - Two Factor Authentication -

func (app *App) resendTwoFactorAuthHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// Get payload data
	var payload payloads.SendVerificationReq

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := app.dataValidator.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Request pasword reset
	if err := app.service.Auth.SendVerification(ctx, auth.TwoFactorAuth, payload.Email); err != nil {
		app.parseError(w, r, err)
		return
	}

	//* No content
	w.WriteHeader(http.StatusNoContent)
}
