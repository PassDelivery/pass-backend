package restaurants_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (s *RestaurantsService) GetRestaurant(
	ctx context.Context,
	restID uuid.UUID,
) (domain.RestaurantCard, error) {
	rest, err := s.restaurantsRepository.GetRestaurant(ctx, restID)
	if err != nil {
		return domain.RestaurantCard{}, fmt.Errorf("get restaurant in service: %w", err)
	}

	return rest, nil
}
