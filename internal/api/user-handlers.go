package api

import (
	"errors"
	"net/http"

	jsonutils "github.com/FelipeBelloDultra/go-bid/internal/json-utils"
	"github.com/FelipeBelloDultra/go-bid/internal/services"
	"github.com/FelipeBelloDultra/go-bid/internal/use-case/user"
)

func (api *API) handleSignUpUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJSON[user.CreateUserReq](r)
	if err != nil {
		_ = jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.CreateUser(
		r.Context(),
		data.UserName,
		data.Email,
		data.Password,
		data.Bio,
	)
	if err != nil {
		if errors.Is(err, services.ErrDuplicatedEmailOrPassword) {
			_ = jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error": "email or username already exists",
			})
		}
	}

	_ = jsonutils.EncodeJSON(w, r, http.StatusCreated, map[string]any{
		"id": id,
	})
}

func (api *API) handleSignInUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}

func (api *API) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	panic("not implemented")
}
