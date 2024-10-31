package api

import "github.com/go-chi/chi/v5"

func (api *API) BindRoutes() {
	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Post("/sign-up", api.handleSignUpUser)
				r.Post("/sign-in", api.handleSignInUser)
				r.Post("/logout", api.handleLogoutUser)
			})
		})
	})
}
