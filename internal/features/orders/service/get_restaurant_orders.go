package orders_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *OrdersService) GetRestaurantOrders(
	ctx context.Context,
	restID uuid.UUID,
) ([]domain.OrderResponse, error) {
	ownerID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get id from context: %w", err)
	}

	isOwner, err := s.ordersRepository.IsRestaurantOwner(ctx, restID, ownerID)
	if err != nil {
		return nil, fmt.Errorf("check restaurant owner: %w", err)
	}
	if !isOwner {
		return nil, fmt.Errorf("restaurant not found: %w", core_errors.ErrNotFound)
	}

	orders, err := s.ordersRepository.GetRestaurantOrders(ctx, restID)
	if err != nil {
		return nil, fmt.Errorf("get restaurant orders in service: %w", err)
	}

	return orders, nil
}
