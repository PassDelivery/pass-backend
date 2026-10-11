package restaurants_transport_http

import (
	"net/http"

	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
)

func (h *RestaurantsHTTPHandler) GetRestaurants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke GetRestaurants handler")

	rests, err := h.restaurantsService.GetRestaurants(ctx)
	if err != nil {
		responseWriter.ErrorResponse(err, "get restaurants")
		return
	}

	responseWriter.JSONResponse(restaurantsToDTO(rests), http.StatusOK)
}
