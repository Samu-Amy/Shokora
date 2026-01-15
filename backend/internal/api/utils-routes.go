package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// - Params -

func (app *App) getIdFromParam(r *http.Request, idParamName string) (int64, error) {
	idParam := chi.URLParam(r, idParamName)

	resourceId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return -1, err
	}

	return resourceId, nil
}
