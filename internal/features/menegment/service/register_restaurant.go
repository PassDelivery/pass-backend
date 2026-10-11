package menegment_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
)

func (s *MenegmentService) RegisterRestaurant(
	ctx context.Context,
	restaurant domain.RestaurantRequest,
) (domain.RestaurantResponse, error) {
	id, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.RestaurantResponse{}, fmt.Errorf("get id from context: %w", err)
	}

	restaurantResp, err := s.menegmentRepository.RegisterRestaurant(ctx, restaurant, id)
	if err != nil {
		return domain.RestaurantResponse{}, fmt.Errorf("register restaurant in service: %w", err)
	}

	return restaurantResp, nil
}
