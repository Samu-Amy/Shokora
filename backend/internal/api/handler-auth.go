package api

import (
	"net/http"
	"time"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	"github.com/Samu-Amy/Shokora/internal/auth"
	"github.com/Samu-Amy/Shokora/internal/errorcodes"
	"github.com/Samu-Amy/Shokora/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ----- REGISTER -----

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
	var payload payloads.RegisterUserReqPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Hash password
	hashedPassword, err := app.hashPassword(payload.Password)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// Create user from payload data
	user := &store.User{
		FirstName:    payload.FirstName,
		LastName:     payload.LastName,
		Email:        payload.Email,
		PasswordHash: hashedPassword,
		ImageUrl:     payload.ImageUrl,
		BirthDate:    payload.BirthDate,
	}

	// Create Response Payload
	resPayload := payloads.RegisterUserResPayload{}

	// Create user in db
	if err := app.store.User.Create(ctx, user); err != nil { // TODO: fare transaction per creazione user, stats and settings (oppure crearle nell'update se non esistono)
		app.parseError(w, r, err)
		return
	}

	resPayload.User = payloads.CreateUserResPayload(user) // Add user to payload

	// TODO: genera token auth e setta i cookie (basta user.Id)

	// Generate verificationTokens (Magic Link and OTP)
	verificationTokens, err := app.tokenAuthenticator.CreateVerificationTokens(auth.EmailVerification)
	if err != nil {
		app.logger.Warnw("error generating verification tokens", "error", err)

		resPayload.Error = errorcodes.ErrVerification.Error() // Add error to payload

		//* Return user, verificationID and error
		if err := app.jsonResponse(w, http.StatusCreated, resPayload); err != nil {
			app.internalServerError(w, r, err)
		}
		return
	}

	// Create Email Verification Tokens
	verificationId, err := app.service.Auth.CreateVerificationTokensWithRetries(ctx, user, verificationTokens)
	if err != nil {
		app.logger.Warnw("error creating email verification tokens in db", "error", err)

		resPayload.Error = errorcodes.ErrVerification.Error() // Add error to payload

		//* Return user, verificationID and error
		if err := app.jsonResponse(w, http.StatusCreated, resPayload); err != nil {
			app.internalServerError(w, r, err)
		}
		return
	}

	resPayload.VerificationId = verificationId // Ad verification id to payload

	// Send email
	err = app.SendVerificationEmail(
		ctx,
		auth.EmailVerification,
		user.FirstName,
		user.Email,
		verificationTokens.PlainMagicLinkToken,
		verificationTokens.PlainOTP,
		verificationTokens.MagicLinkTokenExp,
		verificationTokens.OTPExp,
	)
	if err != nil {
		app.logger.Warnw("error sending welcome email", "error", err)

		resPayload.Error = errorcodes.ErrEmailNotSent.Error() // Add error to payload

		//* Return user, verificationID and error
		if err := app.jsonResponse(w, http.StatusCreated, resPayload); err != nil {
			app.internalServerError(w, r, err)
		}
		return

		// TODO: dire di riprovare più tardi? -> l'utente può accedere ma non può ordinare (ha come opzioni di re-inviare la mail di verifica oppure eliminare l'account (e il token))
	}

	app.logger.Info("User and Tokens created, Email sent successfully")

	// TODO: ricordati di scrivere di controllare nello spam (aggiungere timer al tasto per reinviare la mail (?))

	// TODO: setta cookie auth (prima di ogni return)

	//* Return user and verificationID
	if err := app.jsonResponse(w, http.StatusCreated, resPayload); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// ----- LOGIN -----

func (app *App) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: gestisci casi user non verificato e user verificato ma con 2fa richiesta (se 2fa -> "verify-2fa[/{token}]" -> generate auth tokens ("tokens"), se no 2fa -> generate auth tokens ("tokens"))

	// TODO: ritorna user (se non verificato ritorna RegisterUserResPayload?)
}

// ----- EMAIL VERIFICATION -----

func (app *App) verifyEmailWithTokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get "token" param
	token := chi.URLParam(r, verificationTokenParam)

	// Hash token
	hashedToken := app.tokenAuthenticator.HashMagicLinkToken(&token)

	// Verify
	if err := app.service.Auth.VerifyEmailWithToken(ctx, hashedToken); err != nil {
		app.logger.Warnw("Error with Email Verification using Token", "error", err)

		app.parseError(w, r, err) // TODO: nel FRONTEND dire che "non è valido o è scaduto" (non specificare quale dei due)
		return
	}

	//* No content
	w.WriteHeader(http.StatusNoContent)
}

func (app *App) verifyEmailWithOTPHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload data
	var payload payloads.OTPVerificationReqPayload // TODO: nel FRONTEND ricorda di inviare l'otp come stringa

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if len(payload.OTP) != int(app.config.Auth.OTP.Length) {
		app.badRequestError(w, r, errorcodes.ErrInvalid) // Invalid token
		return
	}

	// Hash OTP
	hashedOTP := app.tokenAuthenticator.HashOTP(payload.OTP, auth.EmailVerification)

	// Verify
	if err := app.service.Auth.VerifyEmailWithOTP(ctx, payload.VerificationId, hashedOTP, app.config.Auth.OTP.MaxAttempts); err != nil {
		app.logger.Warnw("Error with Email Verification using OTP", "error", err)

		app.parseError(w, r, err) // TODO: nel FRONTEND dire che "non è valido o è scaduto" (non specificare quale dei due)
		return
	}

	//* No content
	w.WriteHeader(http.StatusNoContent)
}

// ----- TOKENS -----

// TODO: da "incorporare" in register, login (se no 2fa) e in 2fa (se attiva)
func (app *App) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload data
	var payload payloads.LoginUserReqPayload

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
		app.parseError(w, r, err) // TODO: FRONTEND - non dire se l'email esiste o meno
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(payload.Password))
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

	token, err := app.jwtAuthenticator.GenerateJWTToken(claims)
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
