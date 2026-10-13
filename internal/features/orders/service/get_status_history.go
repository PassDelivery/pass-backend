package orders_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *OrdersService) GetStatusHistory(
	ctx context.Context,
	orderID uuid.UUID,
) ([]domain.OrderStatusChange, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get id from context: %w", err)
	}

	order, err := s.ordersRepository.GetOrder(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}

	if order.UserID != userID {
		return nil, fmt.Errorf("order not found: %w", core_errors.ErrNotFound)
	}

	history, err := s.ordersRepository.GetStatusHistory(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("get status history in service: %w", err)
	}

	return history, nil
}
