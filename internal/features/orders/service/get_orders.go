package orders_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
)

func (s *OrdersService) GetOrders(
	ctx context.Context,
) ([]domain.OrderResponse, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get id from context: %w", err)
	}

	orders, err := s.ordersRepository.GetUserOrders(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get orders in service: %w", err)
	}

	return orders, nil
}
