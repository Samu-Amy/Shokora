package api

import (
	"net/http"
)

func (app *App) CheckHealth(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
