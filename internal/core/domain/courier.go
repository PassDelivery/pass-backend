package domain

import (
	"time"

	"github.com/google/uuid"
)

type CourierStatus string

const (
	CourierStatusOffline CourierStatus = "offline"
	CourierStatusFree    CourierStatus = "free"
	CourierStatusBusy    CourierStatus = "busy"
)

func (s CourierStatus) CanBeSetByCourier() bool {
	return s == CourierStatusOffline || s == CourierStatusFree
}

type DeliveryStatus string

const (
	DeliveryStatusOffered   DeliveryStatus = "offered"
	DeliveryStatusAccepted  DeliveryStatus = "accepted"
	DeliveryStatusRefused   DeliveryStatus = "refused"
	DeliveryStatusPickedUp  DeliveryStatus = "picked_up"
	DeliveryStatusDelivered DeliveryStatus = "delivered"
)

var deliveryTransitions = map[DeliveryStatus]DeliveryStatus{
	DeliveryStatusAccepted: DeliveryStatusPickedUp,
	DeliveryStatusPickedUp: DeliveryStatusDelivered,
}

var deliveryOrderStatus = map[DeliveryStatus]OrderStatus{
	DeliveryStatusPickedUp:  OrderStatusDelivering,
	DeliveryStatusDelivered: OrderStatusDelivered,
}

func (s DeliveryStatus) CanCourierMoveTo(next DeliveryStatus) bool {
	return deliveryTransitions[s] == next
}

func (s DeliveryStatus) OrderStatus() OrderStatus {
	return deliveryOrderStatus[s]
}

type Courier struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	TransportType string
	Status        CourierStatus
}

type Delivery struct {
	ID              uuid.UUID
	Status          DeliveryStatus
	OfferedAt       time.Time
	RespondedAt     *time.Time
	DeliveredAt     *time.Time
	OrderID         uuid.UUID
	OrderStatus     OrderStatus
	RestName        string
	RestAddress     string
	DeliveryAddress string
	Comment         *string
	CustomerName    string
	CustomerPhone   *string
	TotalKopecks    int64
	Items           []OrderItem
}
