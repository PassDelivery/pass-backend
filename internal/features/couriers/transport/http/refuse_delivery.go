package couriers_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *CouriersHTTPHandler) RefuseDelivery(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke RefuseDelivery handler")

	id := r.PathValue("id")
	deliveryID, err := uuid.Parse(id)
	if err != nil {
		responseWriter.ErrorResponse(fmt.Errorf("parse delivery id: %w", core_errors.ErrInvalidArgument), "parse to uuid")
		return
	}

	if err := h.couriersService.RefuseDelivery(ctx, deliveryID); err != nil {
		responseWriter.ErrorResponse(err, "refuse delivery")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
