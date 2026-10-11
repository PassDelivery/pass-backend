package menegment_transport_http

import (
	"net/http"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_request "github.com/PassDelivery/pass-backend/internal/core/transport/http/request"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *MenegmentHTTPHandler) PatchRestaurant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke PatchRestaurant handler")

	id := r.PathValue("id")
	restID, err := uuid.Parse(id)
	if err != nil {
		responseWriter.ErrorResponse(err, "parse to uuid")
		return
	}

	var rest RestaurantPatchDTO

	if err := core_http_request.DecodeAndValidateRequest(r, &rest); err != nil {
		responseWriter.ErrorResponse(err, "decode or validate")
		return
	}


	restaurant, err := h.menegmentService.PatchRestaurant(ctx, domaiRestPatchFromDTO(rest), restID)
	if err != nil {
		responseWriter.ErrorResponse(err, "patch restaurant")
		return
	}

	responseWriter.JSONResponse(DTOFromDomain(restaurant), http.StatusOK)
}

func domaiRestPatchFromDTO(rest RestaurantPatchDTO) domain.RestaurantPatch {
	return domain.RestaurantPatch{
		Name: rest.Name.ToDomain(),
		Description: rest.Description.ToDomain(),
		Address: rest.Address.ToDomain(),
		Phone: rest.Phone.ToDomain(),
		IsOpen: rest.IsOpen.ToDomain(),
	}
}
