package menegment_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *MenegmentHTTPHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke GetCategories handler")

	id := r.PathValue("id")
	restID, err := uuid.Parse(id)
	if err != nil {
		responseWriter.ErrorResponse(fmt.Errorf("parse restaurant id: %w", core_errors.ErrInvalidArgument), "parse to uuid")
		return
	}

	categories, err := h.menegmentService.GetCategories(ctx, restID)
	if err != nil {
		responseWriter.ErrorResponse(err, "get categories")
		return
	}

	response := make([]CategoryResponseDTO, len(categories))
	for k, v := range categories {
		response[k] = categoryToDTO(v)
	}

	responseWriter.JSONResponse(response, http.StatusOK)
}
