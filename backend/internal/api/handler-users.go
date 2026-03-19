package api

import (
	"context"
	"log"
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
	"github.com/Samu-Amy/Shokora/internal/store/user"
)

// TODO: poter modificare dati utente (soprattutto moter mettere nome/cognome dopo la registrazione con google per poter fixarli nel caso non fossero giusti nell'account google)

// ----- GET -----

func (app *App) getUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user id
	userId, err := app.getInt64FromParam(r, userIdParam)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Get user
	user, err := app.getUserById(ctx, userId)
	if err != nil {
		app.parseError(w, r, err)
		return
	}

	//* Return user
	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// ----- UPDATE -----

func (app *App) updateUserDataHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implementa (per modificare i dati dell'utente -
	// TODO: attenzione a evitare modifica data di nascita per avere offerte più volte all'anno (magari durante la registrazione si sottolinea che non è modificabile) -
	// TODO: modificabile solo per chi non la ha (non messa in registrazione o registrato con OAuth e per qualche motivo manca))

	log.Print("\n\nUpdate User...\n\n")
}

func (app *App) updatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get payload data
	var payload payloads.UpdatePasswordReq

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Validate
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if payloads.IsCommonPassword(payload.NewPassword) {
		app.badRequestError(w, r, domerrors.ErrCommonPassword)
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

// ----- UTILS -----

func (app *App) getUserById(ctx context.Context, userId int64) (*user.User, error) {
	user, err := app.service.User.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
