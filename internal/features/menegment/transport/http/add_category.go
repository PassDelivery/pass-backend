package menegment_transport_http

import (
	"net/http"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_request "github.com/PassDelivery/pass-backend/internal/core/transport/http/request"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *MenegmentHTTPHandler) AddCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke AddCategory handler")

	id := r.PathValue("id")
	restID, err := uuid.Parse(id)
	if err != nil {
		responseWriter.ErrorResponse(err, "sfailed scan")
		return
	}

	var catReq CategoryRequestDTO

	if err := core_http_request.DecodeAndValidateRequest(r, &catReq); err != nil {
		responseWriter.ErrorResponse(err, "decode or validate")
		return
	}

	catResp, err := h.menegmentService.AddCategory(ctx, categoryToDomain(catReq), restID)
	if err != nil {
		responseWriter.ErrorResponse(err, "add category")
		return
	}

	responseWriter.JSONResponse(categoryToDTO(catResp), http.StatusCreated)
}

func categoryToDomain(cat CategoryRequestDTO) domain.CategoryRequest {
	return domain.CategoryRequest{
		Name: cat.Name,
		SortOrder: cat.SortOrder,
	}
}

func categoryToDTO(cat domain.CategoryResponse) CategoryResponseDTO {
	return CategoryResponseDTO{
		ID: cat.ID,
		RestID: cat.RestID,
		Name: cat.Name,
		SortOrder: cat.SortOrder,
	}
}