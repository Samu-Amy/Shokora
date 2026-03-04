package api

import (
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
)

// ----- REGISTER -----

// TODO: l'accesso con google (fai handler apposta) sostituisce solo la parte di autenticazione (login e register (in questo caso fornisce già la verifica della mail, settata a true)), poi la gestione di accesso e sessione è gestita dal mio sistema (?)

/*
Return
  - RegisterUserResPayload

Errors
  - ErrDuplicateEmail
  - ErrMaxRetriesExceeded (magic link verification token generation)
  - ErrEmailNotSent	(magic link verification token)
  - ErrInvalidEmailVars (magic link verification token missing)

Payload Error types
  - ErrVerification
  - ErrEmailNotSent
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
	// ctx := r.Context()
	// TODO: gestisci casi user non verificato e user verificato ma con 2fa richiesta (se 2fa -> "verify-2fa[/{token}]" -> generate auth tokens ("tokens"), se no 2fa -> generate auth tokens ("tokens"))

	// TODO: ritorna user (se non verificato ritorna RegisterUserResPayload?)

	// TODO: se login fare pulizia (eliminare token di sessioni scadute - attenzione agli expires aggiornati (vecchi token scaduti ma nuovi no -> sessione ancora valida), controlla per tutta la sessione)?
	// refreshToken, err := app.generateRefreshToken(ctx, user.Id)
	// if err != nil {
	// 	app.internalServerError(w, r, err)
	// 	return
	// }
}

// ----- EMAIL VERIFICATION -----

func (app *App) verifyEmailWithTokenHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	// // Get "token" param
	// token := chi.URLParam(r, verificationTokenParam)

	// // Hash token
	// hashedToken := auth.HashBase64Token(&token)

	// // Verify
	// if err := app.service.Auth.VerifyEmailWithToken(ctx, hashedToken); err != nil {
	// 	app.logger.Warnw("Error with Email Verification using Token", "error", err)

	// 	app.parseError(w, r, err) // TODO: nel FRONTEND dire che "non è valido o è scaduto" (non specificare quale dei due)
	// 	return
	// }

	// //* No content
	// w.WriteHeader(http.StatusNoContent)
}

func (app *App) verifyEmailWithOTPHandler(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()

	// Get payload data
	var payload payloads.OTPVerificationReq // TODO: nel FRONTEND ricorda di inviare l'otp come stringa

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// if len(payload.OTP) != int(app.config.Auth.OTP.Length) {
	// 	app.badRequestError(w, r, domerrors.ErrInvalid) // Invalid token
	// 	return
	// }

	// // Hash OTP
	// hashedOTP := app.tokenAuthenticator.HashOTP(payload.OTP, auth.EmailVerification)

	// // Verify
	// if err := app.service.Auth.VerifyEmailWithOTP(ctx, payload.VerificationId, hashedOTP, app.config.Auth.OTP.MaxAttempts); err != nil {
	// 	app.logger.Warnw("Error with Email Verification using OTP", "error", err)

	// 	app.parseError(w, r, err) // TODO: nel FRONTEND dire che "non è valido o è scaduto" (non specificare quale dei due)
	// 	return
	// }

	// //* No content
	// w.WriteHeader(http.StatusNoContent)
}

// ----- TOKENS -----

// TODO: da "incorporare" in register, login (se no 2fa) e in 2fa (se attiva)
func (app *App) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: usa service

	// ctx := r.Context()

	// // Get payload data
	// var payload payloads.LoginUserReq

	// if err := readJSON(w, r, &payload); err != nil {
	// 	app.badRequestError(w, r, err)
	// 	return
	// }

	// // Validate
	// if err := Validate.Struct(payload); err != nil {
	// 	app.badRequestError(w, r, err)
	// 	return
	// }

	// // Fetch the user (check if the user exist)
	// user, err := app.service.User.GetByEmail(ctx, payload.Email)
	// if err != nil {
	// 	app.parseError(w, r, err) // TODO: FRONTEND - non dire se l'email esiste o meno
	// 	return
	// }

	// // Compare password
	// err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(payload.Password))
	// if err != nil {
	// 	app.unauthorizedError(w, r, err)
	// 	return
	// }

	// // Generate tokens (and add claims)
	// claims := jwt.MapClaims{
	// 	"sub": user.Id, // subject
	// 	"exp": time.Now().Add(app.config.Auth.Token.AccessTokenExp).Unix(),
	// 	"iat": time.Now().Unix(),              // issued at
	// 	"nbf": time.Now().Unix(),              // not before time
	// 	"iss": app.config.Auth.Token.Issuer,   // issuer
	// 	"aud": app.config.Auth.Token.Audience, // audience
	// }

	// token, err := app.jwtAuthenticator.GenerateJWTToken(claims)
	// if err != nil {
	// 	app.internalServerError(w, r, err)
	// 	return
	// }

	// // TODO: setta cookie invece che inviarlo come payload

	// //* Send token to the client
	// if err := app.jsonResponse(w, http.StatusCreated, token); err != nil {
	// 	app.internalServerError(w, r, err)
	// }
}

// func (app *App) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {

// TODO: opzioni per cookie (da verificare)

// HttpOnly: true
// Secure: true
// SameSite: Strict (refresh) / Lax (access)
// Path: /auth/refresh (per refresh token)
// }
