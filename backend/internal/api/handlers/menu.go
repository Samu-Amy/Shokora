package handlers

import (
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/store"
)

func GetAllMenuProducts(store *store.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(`{ "status": "ok" }`))
		// TODO: sostituisci con writeJSON
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
