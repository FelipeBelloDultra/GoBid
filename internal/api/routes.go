package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api *API) BindRoutes() {
	api.Router.Use(
		middleware.RequestID,
		middleware.Logger,
		middleware.Recoverer,
		api.Sessions.LoadAndSave,
	)

	// csrfMiddleware := csrf.Protect(
	// 	[]byte(os.Getenv("GOBID_CSRF_KEY")),
	// 	csrf.Secure(false), // Dev only
	// )

	// api.Router.Use(csrfMiddleware)

	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// r.Get("/csrf-token", api.HandleGetCSRFToken)
			r.Route("/users", func(r chi.Router) {
				r.Post("/sign-up", api.handleSignUpUser)
				r.Post("/sign-in", api.handleSignInUser)
				r.With(api.AuthMiddleware).Post("/logout", api.handleLogoutUser)
			})
		})
	})
}
