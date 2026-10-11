package menegment_service

import (
	"context"
	"fmt"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *MenegmentService) HideDish(
	ctx context.Context,
	dishID uuid.UUID,
) error {
	ownerID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get id from context: %w", err)
	}

	dish, err := s.menegmentRepository.GetDish(ctx, dishID)
	if err != nil {
		return fmt.Errorf("get dish: %w", err)
	}

	isOwner, err := s.menegmentRepository.IsRestaurantOwner(ctx, dish.RestID, ownerID)
	if err != nil {
		return fmt.Errorf("check restaurant owner: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("dish not found: %w", core_errors.ErrNotFound)
	}

	inOrders, err := s.menegmentRepository.IsDishInActiveOrders(ctx, dishID)
	if err != nil {
		return fmt.Errorf("check active orders: %w", err)
	}
	if inOrders {
		return fmt.Errorf("dish is in active orders: %w", core_errors.ErrConflict)
	}

	if err := s.menegmentRepository.HideDish(ctx, dishID); err != nil {
		return fmt.Errorf("hide dish in service: %w", err)
	}

	return nil
}
