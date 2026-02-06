package api

import (
	"expvar"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *App) initRouter() *chi.Mux {
	router := chi.NewRouter()

	//* - Middlewares - *

	// Generic middlewares
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP) // per questo bisogna configurare bene nginx
	/* (tipo:
		set_real_ip_from 127.0.0.1;
		set_real_ip_from <IP_DEL_LOAD_BALANCER>;
		real_ip_header X-Forwarded-For;
		real_ip_recursive on;
	) */
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// router.Use(middleware.Timeout(60 * time.Second)) // Timeout // TODO: bisogna controllarlo negli handler per evitare panic (eventualmente usarlo solo su alcuni?)

	// TODO: setta comunque (con il dominio)
	// CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: app.config.AllowedOriginsURLs,

		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},

		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
		},

		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true, // cookies
		MaxAge:           300,
	}))

	router.Use(app.rateLimiterMiddleware)
	// router.Use(httprate.LimitByIP(100, 1*time.Minute)) // Rate Limiter (?) // TODO: Controlla implementazione (?)

	//* - Routes - *

	// v1
	router.Route("/api/v1", func(r chi.Router) {
		// - Public Routes (commons) -

		r.Get("/health", app.checkHealthHandler) //! TODO: togli

		// Products
		r.Get("/menu/products", app.getMenuProductsHandler)
		// r.Get("/shop/products", app.GetShopProducts)
		// r.Get("/featured/products", app.GetFeaturedProducts) // quelli in "vetrina" sul sito

		r.Get("/products/{productId}", app.getProductHandler) // TODO: per ora prende da product invece che da menu (va bene?)

		// - Auth Routes -
		r.Route("/auth", func(r chi.Router) {
			// Auth
			r.Post("/user", app.registerUserHandler)
			// r.Post("/login", ...) // TODO: se 2fa -> "verify-2fa[/{token}]" -> generate auth tokens ("tokens"), se no 2fa -> generate auth tokens ("tokens")
			// r.Post("/logout", ...)

			r.Post("/tokens", app.createTokenHandler)
			// r.Post("/refresh", ...)

			// Verifications // TODO: implementa routes per login, reset password, ecc.
			r.Post("/verify-email/otp", app.verifyEmailWithOTPHandler) // TODO: spostare in Auth-Protected Routes (?)
			r.Post("/verify-email/{token}", app.verifyEmailWithTokenHandler)

			// r.Post("/reset-password/otp", ...) // TODO: versione logged (usa user Id) e versione non logged (la quale richiede l'email per poter verificare l'otp (in questo caso legato a email invece che user Id))
			// r.Post("/reset-password/{token}", ...)

			// r.Post("/verify-2fa/otp", app.verify2FAWithOTPHandler) // TODO: come verificare otp e utente (id o email) prima di fa accedere l'utente?
			// r.Post("/verify-2fa/{token}", app.verify2FAWithTokenHandler)

			r.Group(func(r chi.Router) {
				// TODO: usare (o crearne uno simile) middleware auth per ottenere l'utente (?)
				// r.Get("/me" app.getCurrentAuthUserHandler) // TODO: per ottenere i dati dell'utente se autenticato
			})
		})

		// - Auth-Protected Routes -
		r.Group(func(r chi.Router) {
			r.Use(app.authMiddleware) // Auth Middleware
			// TODO: controllo modifiche -> gli utenti possono modificare solo il proprio profilo (solo le info di base, non ruolo o altro (quelli modificabili solo da admin))

			// - Customers Routes -

			// User Data
			r.Route("/user", func(r chi.Router) {
				// r.Use(app.userDataOwnershipMiddleware) // User Data Ownerhip Middleware

				// TODO: aggiungi metodo per ottenere user (direttamente dal middleware, non dal db) (con middleware -> solo se è lo stesso di quello autenticato (cioè ottiene il suo profilo))
				r.Patch("/", app.updateUserDataHandler)

				r.Group(func(r chi.Router) {
					// TODO: fai handlers per otterene le stats (con achievements), coupons ed altro
					r.Use(app.userVerifiedMiddleware) // User Verified Middleware (+ User Data Ownerhip Middleware)
					// r.Get("/stats", app.getUserStatsHandler)
					// r.Get("/coupons", app.getUserCouponsHandler)
				})
			})

			// Shop Orders
			r.Group(func(r chi.Router) {
				r.Use(app.userVerifiedMiddleware)
				// r.Post("/orders", app.shopOrderHandler)
			})

			// - Employee (and Admin) Routes -
			r.Route("/employee", func(r chi.Router) {
				r.Use(app.employeeMiddleware)
				// TODO: aggiungi middleware per permessi (?)

				r.Get("/health", app.checkHealthHandler)

				// r.Get("/orders", ...)

				// Users
				r.Route("/users", func(r chi.Router) {
					// r.Get("/", app.getUsersHandler) // Get list of users

					r.Route("/{userId}", func(r chi.Router) {
						r.Get("/", app.getUserHandler)
						// TODO: nell'handler delete user richiedi la password come sicurezza per eliminare l'account (solo jwt non è abbastanza sicuro)
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

				// TODO: il ruolo dev non può essere settato dall'app (solo "a mano" nel db da docker)
				// TODO: route per vedere metrics particolari, logs ed altre cose legate al development (?)
				r.Get("/debug/vars", expvar.Handler().ServeHTTP)
			})
		})
	})

	return router
}
