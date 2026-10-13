package couriers_transport_http

import (
	"net/http"

	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
)

func (h *CouriersHTTPHandler) GetRating(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke GetRating handler")

	rating, err := h.couriersService.GetRating(ctx)
	if err != nil {
		responseWriter.ErrorResponse(err, "get courier rating")
		return
	}

	responseWriter.JSONResponse(ratingToDTO(rating), http.StatusOK)
}
