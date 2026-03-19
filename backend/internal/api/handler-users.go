package api

import (
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/api/payloads"
	domerrors "github.com/Samu-Amy/Shokora/internal/errors/dom"
	user_repo "github.com/Samu-Amy/Shokora/internal/store/user"
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
	user, err := app.service.User.GetById(ctx, userId)
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

// - Update User data -

func (app *App) updateUserDataHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implementa (per modificare i dati dell'utente -
	// TODO: attenzione a evitare modifica data di nascita per avere offerte più volte all'anno (magari durante la registrazione si sottolinea che non è modificabile) -
	// TODO: modificabile solo per chi non la ha (non messa in registrazione o registrato con OAuth e per qualche motivo manca))
}

// - Update Password -

func (app *App) updatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user from context (auth middleware)
	user, ok := r.Context().Value(userCtx).(*user_repo.User)
	if !ok {
		app.unauthorizedError(w, r, domerrors.ErrUnauthorized)
		return
	}

	// Get payload data
	var payload payloads.UpdatePasswordReq

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// - Validation -

	// Basic validation
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	// Passwords must be different
	if payload.NewPassword == payload.OldPassword {
		app.badRequestError(w, r, domerrors.ErrSamePassword)
		return
	}

	// Password must not be common
	if payloads.IsCommonPassword(payload.NewPassword) {
		app.badRequestError(w, r, domerrors.ErrCommonPassword)
		return
	}

	// Update password
	if err := app.service.Auth.UpdatePassword(ctx, user.Id, &payload); err != nil {
		app.parseError(w, r, err)
		return
	}

	//* No content
	w.WriteHeader(http.StatusNoContent)
}
