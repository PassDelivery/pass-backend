package orders_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *OrdersService) GetOrder(
	ctx context.Context,
	orderID uuid.UUID,
) (domain.OrderResponse, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.OrderResponse{}, fmt.Errorf("get id from context: %w", err)
	}

	order, err := s.ordersRepository.GetOrder(ctx, orderID)
	if err != nil {
		return domain.OrderResponse{}, fmt.Errorf("get order in service: %w", err)
	}

	if order.UserID != userID {
		return domain.OrderResponse{}, fmt.Errorf("order not found: %w", core_errors.ErrNotFound)
	}

	return order, nil
}
