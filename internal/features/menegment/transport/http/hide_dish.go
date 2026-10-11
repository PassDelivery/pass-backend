package menegment_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *MenegmentHTTPHandler) HideDish(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke HideDish handler")

	id := r.PathValue("id")
	dishID, err := uuid.Parse(id)
	if err != nil {
		responseWriter.ErrorResponse(fmt.Errorf("parse dish id: %w", core_errors.ErrInvalidArgument), "parse to uuid")
		return
	}

	if err := h.menegmentService.HideDish(ctx, dishID); err != nil {
		responseWriter.ErrorResponse(err, "hide dish")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
