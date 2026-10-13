package orders_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *OrdersHTTPHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke CancelOrder handler")

	id := r.PathValue("id")
	orderID, err := uuid.Parse(id)
	if err != nil {
		responseWriter.ErrorResponse(fmt.Errorf("parse order id: %w", core_errors.ErrInvalidArgument), "parse to uuid")
		return
	}

	order, err := h.ordersService.CancelOrder(ctx, orderID)
	if err != nil {
		responseWriter.ErrorResponse(err, "cancel order")
		return
	}

	responseWriter.JSONResponse(orderToDTO(order), http.StatusOK)
}
