package menegment_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *MenegmentService) PatchDish(
	ctx context.Context,
	patch domain.DishPatch,
	dishID uuid.UUID,
) (domain.DishResponse, error) {
	ownerID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.DishResponse{}, fmt.Errorf("get id from context: %w", err)
	}

	dish, err := s.menegmentRepository.GetDish(ctx, dishID)
	if err != nil {
		return domain.DishResponse{}, fmt.Errorf("get dish: %w", err)
	}

	isOwner, err := s.menegmentRepository.IsRestaurantOwner(ctx, dish.RestID, ownerID)
	if err != nil {
		return domain.DishResponse{}, fmt.Errorf("check restaurant owner: %w", err)
	}
	if !isOwner {
		return domain.DishResponse{}, fmt.Errorf("dish not found: %w", core_errors.ErrNotFound)
	}

	dishReq := domain.DishRequest{
		CategoryID:         dish.CategoryID,
		Name:               dish.Name,
		Description:        dish.Description,
		PriceKopecks:       dish.PriceKopecks,
		CookingTimeMinutes: dish.CookingTimeMinutes,
		IsAvailable:        dish.IsAvailable,
	}

	if err := dishReq.ApplyPatch(patch); err != nil {
		return domain.DishResponse{}, fmt.Errorf("apply patch to dish: %w", err)
	}

	if patch.CategoryID.Set {
		inRest, err := s.menegmentRepository.IsCategoryInRestaurant(ctx, dishReq.CategoryID, dish.RestID)
		if err != nil {
			return domain.DishResponse{}, fmt.Errorf("check category: %w", err)
		}
		if !inRest {
			return domain.DishResponse{}, fmt.Errorf("category not in restaurant: %w", core_errors.ErrInvalidArgument)
		}
	}

	if dish.IsAvailable && !dishReq.IsAvailable {
		inOrders, err := s.menegmentRepository.IsDishInActiveOrders(ctx, dishID)
		if err != nil {
			return domain.DishResponse{}, fmt.Errorf("check active orders: %w", err)
		}
		if inOrders {
			return domain.DishResponse{}, fmt.Errorf("dish is in active orders: %w", core_errors.ErrConflict)
		}
	}

	dishResp, err := s.menegmentRepository.PatchDish(ctx, dishReq, dishID)
	if err != nil {
		return domain.DishResponse{}, fmt.Errorf("patch dish in service: %w", err)
	}

	return dishResp, nil
}
