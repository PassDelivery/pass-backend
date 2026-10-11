package orders_transport_http

import (
	"fmt"
	"time"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
)

type OrderStatusRequestDTO struct {
	Status string `json:"status"`
}

func (d *OrderStatusRequestDTO) Validate() error {
	if !domain.OrderStatus(d.Status).IsValid() {
		return fmt.Errorf("unknown status %q: %w", d.Status, core_errors.ErrInvalidArgument)
	}
	return nil
}

type OrderItemDTO struct {
	ID               uuid.UUID `json:"id"`
	DishID           uuid.UUID `json:"dish_id"`
	DishName         string    `json:"dish_name"`
	Quantity         int       `json:"quantity"`
	PriceAtOrderTime int64     `json:"price_at_order_time"`
}

type OrderResponseDTO struct {
	ID              uuid.UUID      `json:"id"`
	UserID          uuid.UUID      `json:"user_id"`
	RestID          uuid.UUID      `json:"restaurant_id"`
	Status          string         `json:"status"`
	TotalKopecks    int64          `json:"total_kopecks"`
	DeliveryAddress string         `json:"delivery_address"`
	Comment         *string        `json:"comment"`
	Items           []OrderItemDTO `json:"items"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func orderToDTO(order domain.OrderResponse) OrderResponseDTO {
	items := make([]OrderItemDTO, len(order.Items))
	for k, v := range order.Items {
		items[k] = OrderItemDTO{
			ID:               v.ID,
			DishID:           v.DishID,
			DishName:         v.DishName,
			Quantity:         v.Quantity,
			PriceAtOrderTime: v.PriceAtOrderTime,
		}
	}

	return OrderResponseDTO{
		ID:              order.ID,
		UserID:          order.UserID,
		RestID:          order.RestID,
		Status:          string(order.Status),
		TotalKopecks:    order.TotalKopecks,
		DeliveryAddress: order.DeliveryAddress,
		Comment:         order.Comment,
		Items:           items,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}
