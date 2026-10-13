package orders_transport_http

import (
	"net/http"

	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_request "github.com/PassDelivery/pass-backend/internal/core/transport/http/request"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
)

func (h *OrdersHTTPHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke CreateOrder handler")

	var orderReq CreateOrderRequestDTO

	if err := core_http_request.DecodeAndValidateRequest(r, &orderReq); err != nil {
		responseWriter.ErrorResponse(err, "decode or validate")
		return
	}

	order, err := h.ordersService.CreateOrder(ctx, createOrderToDomain(orderReq))
	if err != nil {
		responseWriter.ErrorResponse(err, "create order")
		return
	}

	responseWriter.JSONResponse(orderToDTO(order), http.StatusCreated)
}
