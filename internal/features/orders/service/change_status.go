package orders_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *OrdersService) ChangeOrderStatus(
	ctx context.Context,
	orderID uuid.UUID,
	status domain.OrderStatus,
) (domain.OrderResponse, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.OrderResponse{}, fmt.Errorf("get id from context: %w", err)
	}

	current, ownerID, err := s.ordersRepository.GetOrderStatus(ctx, orderID)
	if err != nil {
		return domain.OrderResponse{}, fmt.Errorf("get order status: %w", err)
	}

	if ownerID != userID {
		return domain.OrderResponse{}, fmt.Errorf("order not found: %w", core_errors.ErrNotFound)
	}

	if !current.CanRestaurantMoveTo(status) {
		return domain.OrderResponse{}, fmt.Errorf("cant move order from %s to %s: %w", current, status, core_errors.ErrConflict)
	}

	order, err := s.ordersRepository.ChangeOrderStatus(ctx, orderID, current, status, userID)
	if err != nil {
		return domain.OrderResponse{}, fmt.Errorf("change order status in service: %w", err)
	}

	return order, nil
}
