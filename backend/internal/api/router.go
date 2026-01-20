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
			r.Post("/user", app.registerUserHandler)
			r.Post("/verify-email/{token}", app.verifyEmailHandler)
			r.Post("/token", app.createTokenHandler)
			// TODO: implementa routes per login, reset password, ecc.
			// r.Post("/login", ...)
			// r.Post("/refresh", ...)
			// r.Post("/reset-password", ...)
			// r.Post("/logout", ...)
		})

		// - Auth-Protected Routes -
		r.Group(func(r chi.Router) {
			r.Use(app.authMiddleware)
			// TODO: controllo modifiche -> gli utenti possono modificare solo il proprio profilo (solo le info di base, non ruolo o altro (quelli modificabili solo da admin))

			// - Customers Routes -

			// User Data
			r.Route("/users/{userId}", func(r chi.Router) {
				r.Use(app.updateUserMiddleware)

				r.Patch("/", app.updateUserDataHandler)
			})

			// Menu Orders

			// Shop Orders

			// - Employee (and Admin) Routes -
			r.Route("/employee", func(r chi.Router) {
				r.Use(app.employeeMiddleware)
				// TODO: aggiungi middleware per permessi (?)

				// r.Get("/orders", ...)

				// Users
				r.Route("/users", func(r chi.Router) {
					// r.Get("/", app.getUsersHandler) // Get list of users

					r.Route("/{userId}", func(r chi.Router) {
						r.Get("/", app.getUserHandler)
					})
				})

				// TODO: usare middleware per i permessi (employyes non possono modificare questi dati se non hanno i permessi)
				// Products
				// r.Group(func(r chi.Router) {
				// r.Use(ProductsManagementAuthorization) // middleware per gestione permessi
				r.Route("/products", func(r chi.Router) {
					r.Post("/", app.createProductHandler)

					r.Route("/{productId}", func(r chi.Router) {
						r.Patch("/", app.updateProductHandler)
						r.Delete("/", app.deleteProductHandler)
					})
				})
				// })
			})

			// - Admin Routes -
			r.Route("/admin", func(r chi.Router) {
				r.Use(app.adminMiddleware)

				// r.Get("/users", ...)
				// r.Patch("/users/{userId}", ...) // gestione users (ruoli e permessi)
			})

			// - Developer Routes -
			r.Route("/dev", func(r chi.Router) {
				r.Use(app.devMiddleware)

				// TODO: route per vedere metrics particolari, logs ed altre cose legate al development (?)
			})
		})
	})

	return router
}
