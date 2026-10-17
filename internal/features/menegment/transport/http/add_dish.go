package menegment_transport_http

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

func (h *MenegmentHTTPHandler) AddDish(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke AddDish handler")

	id := r.PathValue("id")
	restID, err := uuid.Parse(id)
	if err != nil {
		responseWriter.ErrorResponse(fmt.Errorf("parse restaurant id: %w", core_errors.ErrInvalidArgument), "parse to uuid")
		return
	}

	var dishReq DishRequestDTO

	if err := core_http_request.DecodeAndValidateRequest(r, &dishReq); err != nil {
		responseWriter.ErrorResponse(err, "decode or validate")
		return
	}

	dish, err := dishToDomain(dishReq)
	if err != nil {
		responseWriter.ErrorResponse(err, "parse category id")
		return
	}

	dishResp, err := h.menegmentService.AddDish(ctx, dish, restID)
	if err != nil {
		responseWriter.ErrorResponse(err, "add dish")
		return
	}

	responseWriter.JSONResponse(dishToDTO(dishResp), http.StatusCreated)
}

func dishToDomain(dish DishRequestDTO) (domain.DishRequest, error) {
	catID, err := uuid.Parse(dish.CategoryID)
	if err != nil {
		return domain.DishRequest{}, fmt.Errorf("parse category id: %w", core_errors.ErrInvalidArgument)
	}

	return domain.DishRequest{
		CategoryID:         catID,
		Name:               dish.Name,
		Description:        dish.Description,
		PriceKopecks:       dish.PriceKopecks,
		CookingTimeMinutes: dish.CookingTimeMinutes,
		IsAvailable:        dish.IsAvailable,
	}, nil
}

func dishToDTO(dish domain.DishResponse) DishResponseDTO {
	return DishResponseDTO{
		ID:                 dish.ID,
		RestID:             dish.RestID,
		CategoryID:         dish.CategoryID,
		Name:               dish.Name,
		Description:        dish.Description,
		PriceKopecks:       dish.PriceKopecks,
		CookingTimeMinutes: dish.CookingTimeMinutes,
		IsAvailable:        dish.IsAvailable,
		PhotoURL:           dish.PhotoURL,
		CreatedAt:          dish.CreatedAt,
	}
}
