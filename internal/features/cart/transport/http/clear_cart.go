package cart_transport_http

import (
	"net/http"

	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
)

func (h *CartHTTPTransport) ClearCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke ClearCart handler")

	if err := h.cartService.ClearCart(ctx); err != nil {
		responseWriter.ErrorResponse(err, "clear cart")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
