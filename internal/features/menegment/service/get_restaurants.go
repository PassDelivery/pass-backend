package menegment_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
)

func (s *MenegmentService) GetRestaurants(
	ctx context.Context,
) ([]domain.RestaurantResponse, error) {
	ownerID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get id from context: %w", err)
	}

	restaurants, err := s.menegmentRepository.GetRestaurants(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("get restaurants in service: %w", err)
	}

	return restaurants, nil
}
