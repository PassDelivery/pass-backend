package orders_service

import (
	"context"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

type OrdersService struct {
	ordersRepository OrdersRepository
}

type OrdersRepository interface {
	GetOrderStatus(
		ctx context.Context,
		orderID uuid.UUID,
	) (domain.OrderStatus, uuid.UUID, error)
	ChangeOrderStatus(
		ctx context.Context,
		orderID uuid.UUID,
		from domain.OrderStatus,
		to domain.OrderStatus,
		changedBy uuid.UUID,
	) (domain.OrderResponse, error)
}

func NewOrdersService(or OrdersRepository) *OrdersService {
	return &OrdersService{
		ordersRepository: or,
	}
}
