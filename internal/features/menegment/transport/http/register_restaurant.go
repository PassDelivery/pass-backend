package menegment_transport_http

import (
	"net/http"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_request "github.com/PassDelivery/pass-backend/internal/core/transport/http/request"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
)

func (h *MenegmentHTTPHandler) RegisterRestaurant(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke RegisterRestaurant handler")

	var restaurantRequest RegisterRestaurantRequestDTO

	if err := core_http_request.DecodeAndValidateRequest(r, &restaurantRequest); err != nil {
		responseWriter.ErrorResponse(err, "decode or validate")
		return
	}

	restaurant, err := h.menegmentService.RegisterRestaurant(ctx, domainFromDTO(restaurantRequest))
	if err != nil {
		responseWriter.ErrorResponse(err, "register restaurant")
		return
	}

	responseWriter.JSONResponse(DTOFromDomain(restaurant), http.StatusCreated)

}

func domainFromDTO(restaurant RegisterRestaurantRequestDTO) domain.RestaurantRequest {
	return domain.RestaurantRequest{
		Name:        restaurant.Name,
		Description: restaurant.Description,
		Address:     restaurant.Address,
		Phone:       restaurant.Phone,
		IsOpen:      restaurant.IsOpen,
	}
}

func DTOFromDomain(restaurant domain.RestaurantResponse) RestaurantResponseDTO {
	return RestaurantResponseDTO{
		Name:        restaurant.Name,
		Description: restaurant.Description,
		Address:     restaurant.Address,
		Phone:       restaurant.Phone,
		IsOpen:      restaurant.IsOpen,
		ID:          restaurant.ID,
		OwnerID:     restaurant.OwnerID,
		CreatedAt:   restaurant.CreatedAt,
	}
}
