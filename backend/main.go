package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

//TODO: JWT in HTTP only cookies (no in local storage per evitare XSS) -> attenzione a CSRF (cross origin requests)

// TODO: fai test con/senza redis (sia con dati in cache che non in cache) calcolando il tempo impiegato (?)

// DB Connection string
// connStr := "user=${DEV_POSTGRES_USER} dbname=${DEV_POSTGRES_DB} password=${DEV_POSTGRES_PASSWORD} host=localhost port=5432 sslmode=disable"

func main() {
	router := chi.NewRouter()

	// - Middleware -

	// Generic middlewares
	router.Use(middleware.RequestID)
	// router.Use(middleware.RealIP) //! per questo bisogna configurare bene nginx (?)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(middleware.Timeout(60 * time.Second)) // Timeout //TODO: ricorda di controllare ctx.Done() per ritornare nel caso di timeout

	// CORS
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
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},

		AllowedHeaders: []string{
			"Content-Type",
		},

		// TODO: continua... (aggiungi altri url, methods, ecc.)
	}))

	// - Routes -

	// v1
	router.Route("/api/v1", func(r chi.Router) {
		// Public Routes (commons)
		r.Get("/", handleRoot)
		r.Get("/menu", getMenu)
		r.Get("/menu/product/{productId}", getMenuProduct)

		// Auth Routes
		r.Route("/auth", func(r chi.Router) {
			// r.Use(AuthMiddleware)
			// Customers Routes

			// Employee Routes
			r.Route("/employee", func(r chi.Router) {

			})

			// Admin Routes
			r.Route("/admin", func(r chi.Router) {

			})
		})
	})

	// - Server Start -
	fmt.Println("Listening on http://localhost:8080")
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
