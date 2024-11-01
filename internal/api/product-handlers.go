package api

import (
	"net/http"

	jsonutils "github.com/FelipeBelloDultra/go-bid/internal/json-utils"
	"github.com/FelipeBelloDultra/go-bid/internal/use-case/product"
	"github.com/google/uuid"
)

func (api *API) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJSON[product.CreateProductReq](r)
	if err != nil {
		_ = jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problems)
		return
	}
	userID, ok := api.Sessions.Get(r.Context(), "AuthenticatedUserId").(uuid.UUID)
	if !ok {
		_ = jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "internal server error",
		})
		return
	}

	id, err := api.ProductService.Create(
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

	jsonutils.EncodeJSON(w, r, http.StatusCreated, map[string]any{
		"product_id": id,
	})
}
