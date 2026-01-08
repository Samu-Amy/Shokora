package api

import (
	"time"

	"github.com/Samu-Amy/Shokora/internal/api/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func initRouter() *chi.Mux {
	router := chi.NewRouter()

	//* - Middlewares - *

	// Generic middlewares
	router.Use(middleware.RequestID)
	// router.Use(middleware.RealIP) //! per questo bisogna configurare bene nginx (?)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(middleware.Timeout(60 * time.Second)) // Timeout //TODO: ricorda di controllare ctx.Done() per ritornare nel caso di timeout
	// router.Use(httprate.LimitByIP(100, 1*time.Minute)) // Rate Limiter (?) // TODO: Controlla implementazione (?)

	// CORS
	// router.Use(cors.Handler(cors.Options{
	// 	AllowedOrigins: []string{
	// 		"http://localhost:5173",
	// 		"http://192.168.0.46:5173",
	// 		"http://localhost:3000",
	// 		"http://192.168.0.46:3000",
	// 	},

	// 	AllowedMethods: []string{
	// 		http.MethodGet,
	// 		http.MethodPost,
	// 		http.MethodPut,
	// 		http.MethodPatch,
	// 		http.MethodDelete,
	// 	},

	// 	AllowedHeaders: []string{
	// 		// "Accept",
	// 		"Authorization",
	// 		"Content-Type",
	// 	},

	// 	// ExposedHeaders:   []string{"Link"},
	// 	// AllowCredentials: false,
	// 	// MaxAge: 300,

	// 	// TODO: continua... (aggiungi altri url, methods, ecc.)
	// }))

	//* - Routes - *

	// v1
	router.Route("/api/v1", func(r chi.Router) {
		// Public Routes (commons)
		r.Get("/", handlers.HandleRoot)
		r.Get("/menu/products", handlers.GetAllMenuProducts)
		r.Get("/menu/products/{productId}", handlers.GetMenuProduct)

		// Auth Routes
		r.Route("/auth", func(r chi.Router) {
			// r.Post("/login", ...)
			// r.Post("/refresh", ...)
			// r.Post("/reset-password", ...)
			// r.Post("/logout", ...)
		})

		// Auth-Protected Routes
		r.Group(func(r chi.Router) {
			// r.Route("/", func(r chi.Router) { //? Usare Group o Route?
			// r.Use(AuthMiddleware)

			// Customers Routes

			// Employee (and Admin) Routes
			r.Route("/employee", func(r chi.Router) {
				// r.Use(EmployeeMiddleware)

				// r.Get("/orders", ...)
			})

			// Admin Routes
			r.Route("/admin", func(r chi.Router) {
				// r.Use(AdminMiddleware)

				// r.Get("/users", ...)
			})
		})
	})

	return router
}
