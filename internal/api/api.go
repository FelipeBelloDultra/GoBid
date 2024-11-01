package api

import (
	"github.com/FelipeBelloDultra/go-bid/internal/services"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

type API struct {
	Router         *chi.Mux
	Sessions       *scs.SessionManager
	UserService    services.UserService
	ProductService services.ProductService
	BidsService    services.BidsService
}
