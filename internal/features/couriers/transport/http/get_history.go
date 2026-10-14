package couriers_transport_http

import (
	"net/http"

	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
)

func (h *CouriersHTTPHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke GetHistory handler")

	deliveries, err := h.couriersService.GetHistory(ctx)
	if err != nil {
		responseWriter.ErrorResponse(err, "get deliveries history")
		return
	}

	responseWriter.JSONResponse(deliveriesToDTO(deliveries), http.StatusOK)
}
