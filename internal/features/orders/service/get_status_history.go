package orders_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
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

	if _, err := s.userOrder(ctx, orderID, userID); err != nil {
		return nil, err
	}

	history, err := s.ordersRepository.GetStatusHistory(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("get status history in service: %w", err)
	}

	return history, nil
}
