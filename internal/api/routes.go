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

	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Post("/sign-up", api.handleSignUpUser)
				r.Post("/sign-in", api.handleSignInUser)
				r.With(api.AuthMiddleware).Post("/logout", api.handleLogoutUser)
			})
		})
	})
}
