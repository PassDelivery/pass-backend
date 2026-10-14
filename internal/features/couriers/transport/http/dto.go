package couriers_transport_http

import (
	"fmt"
	"time"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
)

type CourierStatusRequestDTO struct {
	Status string `json:"status"`
}

func (d *CourierStatusRequestDTO) Validate() error {
	if !domain.CourierStatus(d.Status).CanBeSetByCourier() {
		return fmt.Errorf("status must be offline or free: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

type DeliveryStatusRequestDTO struct {
	Status string `json:"status"`
}

func (d *DeliveryStatusRequestDTO) Validate() error {
	switch domain.DeliveryStatus(d.Status) {
	case domain.DeliveryStatusPickedUp, domain.DeliveryStatusDelivered:
		return nil
	}
	return fmt.Errorf("status must be picked_up or delivered: %w", core_errors.ErrInvalidArgument)
}

type CourierResponseDTO struct {
	Status        string `json:"status"`
	TransportType string `json:"transport_type,omitempty"`
}

type DeliveryItemDTO struct {
	DishName string `json:"dish_name"`
	Quantity int    `json:"quantity"`
}

type DeliveryResponseDTO struct {
	ID              uuid.UUID         `json:"id"`
	Status          string            `json:"status"`
	OfferedAt       time.Time         `json:"offered_at"`
	RespondedAt     *time.Time        `json:"responded_at"`
	DeliveredAt     *time.Time        `json:"delivered_at"`
	OrderID         uuid.UUID         `json:"order_id"`
	OrderStatus     string            `json:"order_status"`
	RestName        string            `json:"restaurant_name"`
	RestAddress     string            `json:"restaurant_address"`
	DeliveryAddress string            `json:"delivery_address"`
	Comment         *string           `json:"comment"`
	CustomerName    string            `json:"customer_name,omitempty"`
	CustomerPhone   *string           `json:"customer_phone,omitempty"`
	TotalKopecks    int64             `json:"total_kopecks"`
	PayoutKopecks   int64             `json:"payout_kopecks"`
	Items           []DeliveryItemDTO `json:"items"`
}

type CourierRatingResponseDTO struct {
	Rating *float64 `json:"rating"`
	Count  int      `json:"count"`
}

func courierToDTO(courier domain.Courier) CourierResponseDTO {
	return CourierResponseDTO{
		Status:        string(courier.Status),
		TransportType: courier.TransportType,
	}
}

func ratingToDTO(rating domain.CourierRating) CourierRatingResponseDTO {
	return CourierRatingResponseDTO{
		Rating: rating.Rating,
		Count:  rating.Count,
	}
}

func deliveryToDTO(d domain.Delivery) DeliveryResponseDTO {
	items := make([]DeliveryItemDTO, len(d.Items))
	for k, v := range d.Items {
		items[k] = DeliveryItemDTO{
			DishName: v.DishName,
			Quantity: v.Quantity,
		}
	}

	return DeliveryResponseDTO{
		ID:              d.ID,
		Status:          string(d.Status),
		OfferedAt:       d.OfferedAt,
		RespondedAt:     d.RespondedAt,
		DeliveredAt:     d.DeliveredAt,
		OrderID:         d.OrderID,
		OrderStatus:     string(d.OrderStatus),
		RestName:        d.RestName,
		RestAddress:     d.RestAddress,
		DeliveryAddress: d.DeliveryAddress,
		Comment:         d.Comment,
		CustomerName:    d.CustomerName,
		CustomerPhone:   d.CustomerPhone,
		TotalKopecks:    d.TotalKopecks,
		PayoutKopecks:   domain.CourierPayout(d.TotalKopecks),
		Items:           items,
	}
}

func deliveriesToDTO(deliveries []domain.Delivery) []DeliveryResponseDTO {
	res := make([]DeliveryResponseDTO, len(deliveries))
	for k, v := range deliveries {
		res[k] = deliveryToDTO(v)
	}

	return res
}
