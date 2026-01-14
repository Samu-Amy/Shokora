package api

import (
	"net/http"
)

func (app *App) GetMenuProducts(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}

// func GetMenuProduct(store *store.Storage) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Println("Req")

// 		productId := chi.URLParam(r, "productId")

// 		w.Header().Set("Content-Type", "application/json")

// 		w.Write([]byte(`{ "product":` + productId + ` }`))
// 	}
// }
