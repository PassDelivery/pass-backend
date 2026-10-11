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

func (h *MenegmentHTTPHandler) PatchDish(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke PatchDish handler")

	id := r.PathValue("id")
	dishID, err := uuid.Parse(id)
	if err != nil {
		responseWriter.ErrorResponse(fmt.Errorf("parse dish id: %w", core_errors.ErrInvalidArgument), "parse to uuid")
		return
	}

	var dish DishPatchDTO

	if err := core_http_request.DecodeAndValidateRequest(r, &dish); err != nil {
		responseWriter.ErrorResponse(err, "decode or validate")
		return
	}

	dishResp, err := h.menegmentService.PatchDish(ctx, dishPatchToDomain(dish), dishID)
	if err != nil {
		responseWriter.ErrorResponse(err, "patch dish")
		return
	}

	responseWriter.JSONResponse(dishToDTO(dishResp), http.StatusOK)
}

func dishPatchToDomain(dish DishPatchDTO) domain.DishPatch {
	catID := domain.Nullable[uuid.UUID]{Set: dish.CategoryID.Set}
	if dish.CategoryID.Value != nil {
		id := uuid.MustParse(*dish.CategoryID.Value)
		catID.Value = &id
	}

	return domain.DishPatch{
		CategoryID:         catID,
		Name:               dish.Name.ToDomain(),
		Description:        dish.Description.ToDomain(),
		PriceKopecks:       dish.PriceKopecks.ToDomain(),
		CookingTimeMinutes: dish.CookingTimeMinutes.ToDomain(),
		IsAvailable:        dish.IsAvailable.ToDomain(),
	}
}
