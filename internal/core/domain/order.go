package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "new"
	OrderStatusAccepted   OrderStatus = "accepted"
	OrderStatusCooking    OrderStatus = "cooking"
	OrderStatusReady      OrderStatus = "ready"
	OrderStatusDelivering OrderStatus = "delivering"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRejected   OrderStatus = "rejected"
)

var restaurantTransitions = map[OrderStatus][]OrderStatus{
	OrderStatusNew:      {OrderStatusAccepted, OrderStatusRejected},
	OrderStatusAccepted: {OrderStatusCooking},
	OrderStatusCooking:  {OrderStatusReady},
}

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusNew, OrderStatusAccepted, OrderStatusCooking, OrderStatusReady,
		OrderStatusDelivering, OrderStatusDelivered, OrderStatusCancelled, OrderStatusRejected:
		return true
	}
	return false
}

func (s OrderStatus) CanRestaurantMoveTo(next OrderStatus) bool {
	for _, st := range restaurantTransitions[s] {
		if st == next {
			return true
		}
	}
	return false
}

type OrderItem struct {
	ID               uuid.UUID
	DishID           uuid.UUID
	DishName         string
	Quantity         int
	PriceAtOrderTime int64
}

type OrderResponse struct {
	ID              uuid.UUID
	UserID          uuid.UUID
	RestID          uuid.UUID
	Status          OrderStatus
	TotalKopecks    int64
	DeliveryAddress string
	Comment         *string
	Items           []OrderItem
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
