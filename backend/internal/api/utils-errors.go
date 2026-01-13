package api

import (
	"log"
	"net/http"
)

func (app *App) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *App) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bar request error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *App) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusNotFound, "not found")
}
