package api

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *App) initRouter() *chi.Mux {
	router := chi.NewRouter()

	//* - Middlewares - *

	// Generic middlewares
	router.Use(middleware.RequestID)
	// router.Use(middleware.RealIP) //! per questo bisogna configurare bene nginx (?)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(middleware.Timeout(60 * time.Second)) // Timeout
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
		// - Public Routes (commons) -

		r.Get("/health", app.checkHealthHandler)

		// Products
		r.Get("/menu/products", app.getMenuProductsHandler)
		// r.Get("/shop/products", app.GetShopProducts)
		// r.Get("/featured/products", app.GetFeaturedProducts) // quelli in "vetrina" sul sito

		r.Get("/products/{productId}", app.getProductHandler) // TODO: per ora prende da product invece che da menu (va bene?)

		// - Auth Routes -
		r.Route("/auth", func(r chi.Router) {
			// r.Post("/login", ...)
			// r.Post("/refresh", ...)
			// r.Post("/reset-password", ...)
			// r.Post("/logout", ...)
		})

		// - Auth-Protected Routes -
		r.Group(func(r chi.Router) {
			// r.Use(AuthMiddleware)
			// TODO: controllo modifiche -> gli utenti possono modificare solo il proprio profilo (solo le info di base, non ruolo o altro (quelli modificabili solo da admin))

			// - Customers Routes -

			// Menu Orders

			// Shop Orders

			// - Employee (and Admin) Routes -
			r.Route("/employee", func(r chi.Router) {
				// r.Use(EmployeeMiddleware)
				// TODO: aggiungi middleware per permessi (?)

				// r.Get("/orders", ...)

				// Users
				r.Route("/users/{userId}", func(r chi.Router) {
					r.Get("/", app.getUserHandler)
				})

				// Products
				r.Route("/products", func(r chi.Router) {
					r.Post("/", app.createProductHandler)

					r.Route("/{productId}", func(r chi.Router) {
						r.Patch("/", app.updateProductHandler)
						r.Delete("/", app.deleteProductHandler)
					})
				})
			})

			// - Admin Routes -
			r.Route("/admin", func(r chi.Router) {
				// r.Use(AdminMiddleware)

				// r.Get("/users", ...)
			})
		})
	})

	return router
}
