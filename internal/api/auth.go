package api

import (
	"net/http"

	jsonutils "github.com/FelipeBelloDultra/go-bid/internal/json-utils"
)

func (api *API) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.Sessions.Exists(r.Context(), "AuthenticatedUserId") {
			jsonutils.EncodeJSON(w, r, http.StatusUnauthorized, map[string]any{
				"message": "must be authenticated",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
