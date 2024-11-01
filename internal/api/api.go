package api

import (
	"github.com/FelipeBelloDultra/go-bid/internal/services"
	"github.com/go-chi/chi/v5"
)

type API struct {
	Router      *chi.Mux
	UserService services.UserService
}
