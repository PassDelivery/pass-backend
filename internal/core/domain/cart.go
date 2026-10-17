package domain

import "github.com/google/uuid"

type CartItemRequest struct {
	DishID   uuid.UUID
	Quantity int
}

type CartItem struct {
	ID           uuid.UUID
	DishID       uuid.UUID
	Name         string
	PriceKopecks int64
	Quantity     int
	IsAvailable  bool
	PhotoURL     *string
}

type Cart struct {
	RestID       *uuid.UUID
	RestName     *string
	Items        []CartItem
	TotalKopecks int64
}

func (i CartItem) TotalKopecks() int64 {
	return i.PriceKopecks * int64(i.Quantity)
}
