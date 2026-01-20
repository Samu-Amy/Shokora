package api

import (
	"context"
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/store"
)

// TODO: poter modificare dati utente (soprattutto moter mettere nome/cognome dopo la registrazione con google per poter fixarli nel caso non fossero giusti nell'account google)

var userIdParam = "userId"

// ----- GET -----

func (app *App) getUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get user id
	userId, err := app.getIdFromParam(r, userIdParam)
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

// ----- UTILS -----

func (app *App) getUserById(ctx context.Context, userId int64) (*store.User, error) {
	user, err := app.store.User.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
