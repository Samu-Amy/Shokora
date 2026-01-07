package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetAllMenuProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(`{ "status": "ok" }`))
}

func GetMenuProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Req")

	productId := chi.URLParam(r, "productId")

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(`{ "product":` + productId + ` }`))
}
