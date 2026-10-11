package restaurants_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (s *RestaurantsService) GetMenu(
	ctx context.Context,
	restID uuid.UUID,
) (domain.Menu, error) {
	if _, err := s.restaurantsRepository.GetRestaurant(ctx, restID); err != nil {
		return domain.Menu{}, fmt.Errorf("get restaurant: %w", err)
	}

	cats, err := s.restaurantsRepository.GetCategories(ctx, restID)
	if err != nil {
		return domain.Menu{}, fmt.Errorf("get categories: %w", err)
	}

	dishes, err := s.restaurantsRepository.GetAvailableDishes(ctx, restID)
	if err != nil {
		return domain.Menu{}, fmt.Errorf("get dishes: %w", err)
	}

	byCategory := make(map[uuid.UUID][]domain.MenuDish)
	for _, dish := range dishes {
		byCategory[dish.CategoryID] = append(byCategory[dish.CategoryID], dish)
	}

	menu := domain.Menu{
		RestID:     restID,
		Categories: []domain.MenuCategory{},
	}

	for _, cat := range cats {
		if len(byCategory[cat.ID]) == 0 {
			continue
		}

		cat.Dishes = byCategory[cat.ID]
		menu.Categories = append(menu.Categories, cat)
	}

	return menu, nil
}
