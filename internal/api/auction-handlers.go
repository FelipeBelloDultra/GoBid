package api

import (
	"errors"
	"net/http"

	jsonutils "github.com/FelipeBelloDultra/go-bid/internal/json-utils"
	"github.com/FelipeBelloDultra/go-bid/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *API) handleSubscribeUserToAuction(w http.ResponseWriter, r *http.Request) {
	rawProductId := chi.URLParam(r, "product_id")
	productId, err := uuid.Parse(rawProductId)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusBadRequest, map[string]any{
			"error": "invalid product id",
		})
		return
	}

	_, err = api.ProductService.GetProductByID(r.Context(), productId)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			jsonutils.EncodeJSON(w, r, http.StatusNotFound, map[string]any{
				"error": "product not found",
			})
			return
		}

		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
		return
	}

	userId, ok := api.Sessions.Get(r.Context(), AuthenticationSessionKey).(uuid.UUID)
	if !ok {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
		return
	}

	api.AuctionLobby.Lock()
	room, ok := api.AuctionLobby.Rooms[productId]
	api.AuctionLobby.Unlock()

	if !ok {
		jsonutils.EncodeJSON(w, r, http.StatusBadRequest, map[string]any{
			"error": "auction has ended",
		})
		return
	}

	conn, err := api.WsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
		return
	}

	client := services.NewClient(room, conn, userId)

	room.Register <- client
	// go client.ReadEventLoop()
	// go client.WriteEventLoop()
	for {
	}
}
