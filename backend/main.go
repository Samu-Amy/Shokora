package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

//TODO: JWT in HTTP only cookies (no in local storage per evitare XSS) -> attenzione a CSRF (cross origin requests)

// TODO: fai test con/senza redis (sia con dati in cache che non in cache) calcolando il tempo impiegato

func main() {
	router := chi.NewRouter()

	// TODO: aggiungi middleware generali

	// TODO: aggiungi altri url, methods, ecc.
	router.Use(cors.Handler(cors.Options{

		AllowedOrigins: []string{
			"http://localhost:5173",
			"http://192.168.0.46:5173",
			"http://localhost:3000",
			"http://192.168.0.46:3000",
		},

		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
		},

		AllowedHeaders: []string{
			"Content-Type",
		},

		// TODO: continua...
	}))

	router.Get("/api/v1/", handleRoot)
	router.Get("/api/v1/menu", getMenu)
	router.Get("/api/v1/menu/product/{productId}", getMenuProduct)

	fmt.Println("Listening on http://localhost:8080")
	// err := http.ListenAndServe("0.0.0.0:8080", router) // testing on mobile
	err := http.ListenAndServe(":8080", router)

	if err != nil {
		log.Println(err) // TODO: sistema
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func getMenu(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(`{ "status": "ok" }`))
}

func getMenuProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Req")

	productId := chi.URLParam(r, "productId")

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(`{ "product":` + productId + ` }`))
}
