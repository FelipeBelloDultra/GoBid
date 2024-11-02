package api

import (
	"net/http"

	jsonutils "github.com/FelipeBelloDultra/go-bid/internal/json-utils"
	"github.com/gorilla/csrf"
)

const (
	AuthenticationSessionKey = "AuthenticatedUserId"
)

func (api *API) HandleGetCSRFToken(w http.ResponseWriter, r *http.Request) {
	token := csrf.Token(r)
	jsonutils.EncodeJSON(w, r, http.StatusOK, map[string]string{
		"csrf_token": token,
	})
}

func (api *API) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.Sessions.Exists(r.Context(), AuthenticationSessionKey) {
			jsonutils.EncodeJSON(w, r, http.StatusUnauthorized, map[string]any{
				"message": "must be authenticated",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
