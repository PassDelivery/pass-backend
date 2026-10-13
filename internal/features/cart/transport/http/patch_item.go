package cart_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_request "github.com/PassDelivery/pass-backend/internal/core/transport/http/request"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *CartHTTPTransport) PatchItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke PatchItem handler")

	id := r.PathValue("id")
	itemID, err := uuid.Parse(id)
	if err != nil {
		responseWriter.ErrorResponse(fmt.Errorf("parse item id: %w", core_errors.ErrInvalidArgument), "parse to uuid")
		return
	}

	var itemPatch CartItemPatchDTO

	if err := core_http_request.DecodeAndValidateRequest(r, &itemPatch); err != nil {
		responseWriter.ErrorResponse(err, "decode or validate")
		return
	}

	cart, err := h.cartService.PatchItem(ctx, itemID, itemPatch.Quantity)
	if err != nil {
		responseWriter.ErrorResponse(err, "patch cart item")
		return
	}

	responseWriter.JSONResponse(cartToDTO(cart), http.StatusOK)
}
