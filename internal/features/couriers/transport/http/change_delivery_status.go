package couriers_transport_http

import (
	"fmt"
	"net/http"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_request "github.com/PassDelivery/pass-backend/internal/core/transport/http/request"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *CouriersHTTPHandler) ChangeDeliveryStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke ChangeDeliveryStatus handler")

	id := r.PathValue("id")
	deliveryID, err := uuid.Parse(id)
	if err != nil {
		responseWriter.ErrorResponse(fmt.Errorf("parse delivery id: %w", core_errors.ErrInvalidArgument), "parse to uuid")
		return
	}

	var statusReq DeliveryStatusRequestDTO

	if err := core_http_request.DecodeAndValidateRequest(r, &statusReq); err != nil {
		responseWriter.ErrorResponse(err, "decode or validate")
		return
	}

	delivery, err := h.couriersService.ChangeDeliveryStatus(ctx, deliveryID, domain.DeliveryStatus(statusReq.Status))
	if err != nil {
		responseWriter.ErrorResponse(err, "change delivery status")
		return
	}

	responseWriter.JSONResponse(deliveryToDTO(delivery), http.StatusOK)
}
