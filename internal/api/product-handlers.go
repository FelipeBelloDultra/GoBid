package api

import (
	"context"
	"net/http"

	jsonutils "github.com/FelipeBelloDultra/go-bid/internal/json-utils"
	"github.com/FelipeBelloDultra/go-bid/internal/services"
	"github.com/FelipeBelloDultra/go-bid/internal/use-case/product"
	"github.com/google/uuid"
)

func (api *API) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJSON[product.CreateProductReq](r)
	if err != nil {
		_ = jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problems)
		return
	}
	userID, ok := api.Sessions.Get(r.Context(), AuthenticationSessionKey).(uuid.UUID)
	if !ok {
		_ = jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
		return
	}

	productId, err := api.ProductService.Create(
		r.Context(),
		userID,
		data.ProductName,
		data.Description,
		data.BasePrice,
		data.AuctionEnd,
	)
	if err != nil {
		_ = jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
		return
	}

	ctx, _ := context.WithDeadline(context.Background(), data.AuctionEnd)
	auctionRoom := services.NewAuctionRoom(ctx, productId, api.BidsService)

	go auctionRoom.Run()

	api.AuctionLobby.Lock()
	api.AuctionLobby.Rooms[productId] = auctionRoom
	api.AuctionLobby.Unlock()

	jsonutils.EncodeJSON(w, r, http.StatusCreated, map[string]any{
		"product_id": productId,
		"message":    "acution has started with success",
	})
}
