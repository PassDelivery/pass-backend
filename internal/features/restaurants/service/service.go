package restaurants_service

import (
	"context"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

type RestaurantsService struct {
	restaurantsRepository RestaurantsRepository
}

type RestaurantsRepository interface {
	GetOpenRestaurants(
		ctx context.Context,
	) ([]domain.RestaurantCard, error)
	GetRestaurant(
		ctx context.Context,
		restID uuid.UUID,
	) (domain.RestaurantCard, error)
	GetCategories(
		ctx context.Context,
		restID uuid.UUID,
	) ([]domain.MenuCategory, error)
	GetAvailableDishes(
		ctx context.Context,
		restID uuid.UUID,
	) ([]domain.MenuDish, error)
}

func NewRestaurantsService(rr RestaurantsRepository) *RestaurantsService {
	return &RestaurantsService{
		restaurantsRepository: rr,
	}
}
