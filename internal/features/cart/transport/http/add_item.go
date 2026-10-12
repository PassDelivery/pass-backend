package cart_transport_http

import (
	"net/http"

	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_request "github.com/PassDelivery/pass-backend/internal/core/transport/http/request"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
)

func (h *CartHTTPTransport) AddItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke AddItem handler")

	var itemReq CartItemRequestDTO

	if err := core_http_request.DecodeAndValidateRequest(r, &itemReq); err != nil {
		responseWriter.ErrorResponse(err, "decode or validate")
		return
	}

	cart, err := h.cartService.AddItem(ctx, cartItemToDomain(itemReq))
	if err != nil {
		responseWriter.ErrorResponse(err, "add item to cart")
		return
	}

	responseWriter.JSONResponse(cartToDTO(cart), http.StatusOK)
}
