package api

import (
	"expvar"
	"net/http"

	"github.com/Samu-Amy/Shokora/internal/store/user"
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
	router.Route("/api/v1", func(r chi.Router) { //! Se la route dovesse cambiare, modificare path cookies

		// ----- PUBLIC ROUTES (commons) -----

		r.Get("/health", app.checkHealthHandler) // TODO: togli

		// Products
		r.Get("/menu/products", app.getMenuProductsHandler)
		// r.Get("/shop/products", app.GetShopProducts)
		// r.Get("/featured/products", app.GetFeaturedProducts) // quelli in "vetrina" sul sito

		r.Get("/products/{productId}", app.getProductHandler) // TODO: per ora prende da product invece che da menu (va bene?)

		//

		// ----- AUTH ROUTES -----

		r.Route("/auth", func(r chi.Router) {
			// Auth
			r.Post("/user", app.registerUserHandler)
			r.Post("/login", app.loginUserHandler)
			r.Post("/logout", app.logoutUserHandler)

			// Google
			// r.Post("/google", app.googleRegisterUserHandler)
			// r.Post("/login/google", app.googleLoginUserHandler)

			// Verifications
			// TODO: aggiungere route per il reset della password (con email nel payload, ottiene userId dall'email, crea verification tokens (con anche userId) ed invia email con i token)
			r.Post("/verify-email/otp", app.verifyEmailWithOTPHandler) // TODO: spostare verifiche email in Auth-Protected Routes (?)
			r.Post("/verify-email/{token}", app.verifyEmailWithMagicLinkHandler)
			// r.Post("/verify-email/resend", app.resendEmailVerificationHandler) // TODO: fare così?

			r.Post("/reset-password", app.requestPasswordResetHandler)
			r.Post("/reset-password/otp", app.resetPasswordWithMagicLinkHandler)
			r.Post("/reset-password/{token}", app.resetPasswordWithOTPHandler)

			r.Post("/verify-2fa/otp", app.verifyTwoFactorAuthWithOTPHandler)

			r.Group(func(r chi.Router) {
				// TODO: usare (o crearne uno simile) middleware auth per ottenere l'utente (?)
				// r.Get("/me" app.getCurrentAuthUserHandler) // TODO: per ottenere i dati dell'utente se autenticato
				// r.Get("/settings", ...)
				// r.Patch("/settings", ...)
			})
		})

		//

		// ----- AUTH PROTECTED ROUTES -----

		r.Group(func(r chi.Router) {

			// --- USER ROUTES ---

			r.Use(app.authMiddleware) // Auth Middleware
			// TODO: controllo modifiche -> gli utenti possono modificare solo il proprio profilo (solo le info di base, non ruolo o altro (quelli modificabili solo da admin))

			// - Basic User Routes -

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

			// - Verified User Routes -

			// Shop Orders
			r.Group(func(r chi.Router) {
				r.Use(app.userVerifiedMiddleware)
				// r.Post("/orders", app.createOrderHandler)
			})

			//

			// --- EMPLOYEE ROUTES ---

			r.Route("/employee", func(r chi.Router) {

				// - Basic -

				r.Group(func(r chi.Router) {
					r.Use(app.roleMiddleware(user.RoleEmployee))

					// r.Get("/health", app.checkHealthHandler)

					// r.Get("/orders", ...)

					// Users
					r.Route("/users", func(r chi.Router) {
						// r.Get("/", app.getUsersHandler) // Get list of users

						r.Route("/{userId}", func(r chi.Router) {
							r.Get("/", app.getUserHandler)
							// TODO: nell'handler delete user richiedi la password come sicurezza per eliminare l'account (solo jwt non è abbastanza sicuro)
						})
					})
				})

				// - With Permissions -

				// Products
				r.Route("/products", func(r chi.Router) {
					r.With(app.permissionMiddleware(user.RoleEmployee, user.EmplPermProductAdd)).Post("/", app.createProductHandler)

					r.Route("/{productId}", func(r chi.Router) {
						r.With(app.permissionMiddleware(user.RoleEmployee, user.EmplPermProductUpdate)).Patch("/", app.updateProductHandler)
						r.With(app.permissionMiddleware(user.RoleEmployee, user.EmplPermProductDelete)).Delete("/", app.deleteProductHandler)
					})
				})

				// Stock Products
				r.Route("/stock-products", func(r chi.Router) {
					// r.With(app.permissionMiddleware(user.RoleEmployee, user.EmplPermStockProductAdd)).Post("/", app.createStockProductHandler)

					r.Route("/{productId}", func(r chi.Router) {
						// r.With(app.permissionMiddleware(user.RoleEmployee, user.EmplPermStockProductUpdate)).Patch("/", app.updateStockProductHandler)
						// r.With(app.permissionMiddleware(user.RoleEmployee, user.EmplPermStockProductDelete)).Delete("/", app.deleteStockProductHandler)
					})
				})

				// Orders Products
				r.Route("/orders", func(r chi.Router) {
					// r.With(app.permissionMiddleware(user.RoleEmployee, user.EmplPermOrderGet)).Get("/", app.getOrdersHandler)

					r.Route("/{orderId}", func(r chi.Router) {
						// r.With(app.permissionMiddleware(user.RoleEmployee, user.EmplPermOrderGet)).Get("/", app.getOrderHandler)
						// r.With(app.permissionMiddleware(user.RoleEmployee, user.EmplPermOrderUpdate)).Patch("/", app.updateOrderHandler)
						// r.With(app.permissionMiddleware(user.RoleEmployee, user.EmplPermOrderDelete)).Delete("/", app.deleteOrderHandler)
					})
				})
			})

			//

			// --- ADMIN ROUTES ---

			r.Route("/admin", func(r chi.Router) {
				r.Use(app.roleMiddleware(user.RoleAdmin))

				// r.Get("/users", ...)
				// r.Patch("/users/{userId}", ...) // gestione users (ruoli e permessi)
			})

			//

			// --- DEV ROUTES ---

			r.Route("/dev", func(r chi.Router) {
				r.Use(app.roleMiddleware(user.RoleDev))

				// TODO: il ruolo dev non può essere settato dall'app (solo "a mano" nel db da docker)
				// TODO: route per vedere metrics particolari, logs ed altre cose legate al development (?)
				r.Get("/debug/vars", expvar.Handler().ServeHTTP)
			})
		})
	})

	return router
}
