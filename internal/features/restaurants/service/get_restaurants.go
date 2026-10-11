package restaurants_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
)

func (s *RestaurantsService) GetRestaurants(
	ctx context.Context,
) ([]domain.RestaurantCard, error) {
	rests, err := s.restaurantsRepository.GetOpenRestaurants(ctx)
	if err != nil {
		return nil, fmt.Errorf("get restaurants in service: %w", err)
	}

	return rests, nil
}
