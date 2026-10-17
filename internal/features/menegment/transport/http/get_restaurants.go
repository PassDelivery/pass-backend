package menegment_transport_http

import (
	"net/http"

	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
)

func (h *MenegmentHTTPHandler) GetRestaurants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke GetRestaurants handler")

	restaurants, err := h.menegmentService.GetRestaurants(ctx)
	if err != nil {
		responseWriter.ErrorResponse(err, "get restaurants")
		return
	}

	response := make([]RestaurantResponseDTO, len(restaurants))
	for k, v := range restaurants {
		response[k] = DTOFromDomain(v)
	}

	responseWriter.JSONResponse(response, http.StatusOK)
}
