package restaurants_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *RestaurantsHTTPHandler) GetRestaurant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke GetRestaurant handler")

	id := r.PathValue("id")
	restID, err := uuid.Parse(id)
	if err != nil {
		responseWriter.ErrorResponse(fmt.Errorf("parse restaurant id: %w", core_errors.ErrInvalidArgument), "parse to uuid")
		return
	}

	rest, err := h.restaurantsService.GetRestaurant(ctx, restID)
	if err != nil {
		responseWriter.ErrorResponse(err, "get restaurant")
		return
	}

	responseWriter.JSONResponse(restaurantToDTO(rest), http.StatusOK)
}
