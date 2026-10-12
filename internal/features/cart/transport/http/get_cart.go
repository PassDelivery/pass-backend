package cart_transport_http

import (
	"net/http"

	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
)

func (h *CartHTTPTransport) GetCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke GetCart handler")

	cart, err := h.cartService.GetCart(ctx)
	if err != nil {
		responseWriter.ErrorResponse(err, "get cart")
		return
	}

	responseWriter.JSONResponse(cartToDTO(cart), http.StatusOK)
}
