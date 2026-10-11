package menegment_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *MenegmentService) AddDish(
	ctx context.Context,
	dish domain.DishRequest,
	restID uuid.UUID,
) (domain.DishResponse, error) {
	ownerID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.DishResponse{}, fmt.Errorf("get id from context: %w", err)
	}

	isOwner, err := s.menegmentRepository.IsRestaurantOwner(ctx, restID, ownerID)
	if err != nil {
		return domain.DishResponse{}, fmt.Errorf("check restaurant owner: %w", err)
	}
	if !isOwner {
		return domain.DishResponse{}, fmt.Errorf("restaurant not found: %w", core_errors.ErrNotFound)
	}

	inRest, err := s.menegmentRepository.IsCategoryInRestaurant(ctx, dish.CategoryID, restID)
	if err != nil {
		return domain.DishResponse{}, fmt.Errorf("check category: %w", err)
	}
	if !inRest {
		return domain.DishResponse{}, fmt.Errorf("category not in restaurant: %w", core_errors.ErrInvalidArgument)
	}

	dishResp, err := s.menegmentRepository.AddDish(ctx, dish, restID)
	if err != nil {
		return domain.DishResponse{}, fmt.Errorf("add dish in service: %w", err)
	}

	return dishResp, nil
}
