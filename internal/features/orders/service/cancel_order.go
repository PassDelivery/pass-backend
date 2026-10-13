package orders_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *OrdersService) CancelOrder(
	ctx context.Context,
	orderID uuid.UUID,
) (domain.OrderResponse, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.OrderResponse{}, fmt.Errorf("get id from context: %w", err)
	}

	order, err := s.ordersRepository.GetOrder(ctx, orderID)
	if err != nil {
		return domain.OrderResponse{}, fmt.Errorf("get order: %w", err)
	}

	if order.UserID != userID {
		return domain.OrderResponse{}, fmt.Errorf("order not found: %w", core_errors.ErrNotFound)
	}

	if !order.Status.CanCustomerCancel() {
		return domain.OrderResponse{}, fmt.Errorf("cant cancel order in status %s: %w", order.Status, core_errors.ErrConflict)
	}

	cancelled, err := s.ordersRepository.ChangeOrderStatus(ctx, orderID, order.Status, domain.OrderStatusCancelled, userID)
	if err != nil {
		return domain.OrderResponse{}, fmt.Errorf("cancel order in service: %w", err)
	}

	cancelled.RestName = order.RestName

	return cancelled, nil
}
