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
	CreateOrder(
		ctx context.Context,
		userID uuid.UUID,
		order domain.CreateOrderRequest,
	) (domain.OrderResponse, error)
	GetUserOrders(
		ctx context.Context,
		userID uuid.UUID,
	) ([]domain.OrderResponse, error)
	GetOrder(
		ctx context.Context,
		orderID uuid.UUID,
	) (domain.OrderResponse, error)
	GetStatusHistory(
		ctx context.Context,
		orderID uuid.UUID,
	) ([]domain.OrderStatusChange, error)
	IsRestaurantOwner(
		ctx context.Context,
		restID uuid.UUID,
		ownerID uuid.UUID,
	) (bool, error)
	GetRestaurantOrders(
		ctx context.Context,
		restID uuid.UUID,
	) ([]domain.OrderResponse, error)
}

func NewOrdersService(or OrdersRepository) *OrdersService {
	return &OrdersService{
		ordersRepository: or,
	}
}
