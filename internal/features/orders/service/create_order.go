package orders_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
)

func (s *OrdersService) CreateOrder(
	ctx context.Context,
	order domain.CreateOrderRequest,
) (domain.OrderResponse, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.OrderResponse{}, fmt.Errorf("get id from context: %w", err)
	}

	orderResp, err := s.ordersRepository.CreateOrder(ctx, userID, order)
	if err != nil {
		return domain.OrderResponse{}, fmt.Errorf("create order in service: %w", err)
	}

	return orderResp, nil
}
