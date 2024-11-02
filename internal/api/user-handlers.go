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
		if errors.Is(err, services.ErrDuplicatedEmailOrUsername) {
			_ = jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error": "email or username already exists",
			})
			return
		}
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
		return
	}

	_ = jsonutils.EncodeJSON(w, r, http.StatusCreated, map[string]any{
		"id": id,
	})
}

func (api *API) handleSignInUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJSON[user.LoginUserReq](r)
	if err != nil {
		_ = jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.AuthenticateUser(
		r.Context(),
		data.Email,
		data.Password,
	)

	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			_ = jsonutils.EncodeJSON(w, r, http.StatusUnauthorized, map[string]any{
				"error": "invalid email or password",
			})
			return
		}
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
		return
	}

	err = api.Sessions.RenewToken(r.Context())
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
		return
	}
	api.Sessions.Put(r.Context(), AuthenticationSessionKey, id)

	jsonutils.EncodeJSON(w, r, http.StatusOK, map[string]any{
		"message": "logged in successfully",
	})
}

func (api *API) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	err := api.Sessions.RenewToken(r.Context())
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
		return
	}

	api.Sessions.Remove(r.Context(), AuthenticationSessionKey)
	jsonutils.EncodeJSON(w, r, http.StatusOK, map[string]any{
		"message": "logged out successfully",
	})
}
