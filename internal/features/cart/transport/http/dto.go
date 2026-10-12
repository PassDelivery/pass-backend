package cart_transport_http

import (
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
)

const maxQuantity = 99

type CartItemRequestDTO struct {
	DishID   string `json:"dish_id"`
	Quantity int    `json:"quantity"`
}

func (d *CartItemRequestDTO) Validate() error {
	if _, err := uuid.Parse(d.DishID); err != nil {
		return fmt.Errorf("dish_id is not uuid: %w", core_errors.ErrInvalidArgument)
	}
	if d.Quantity < 1 || d.Quantity > maxQuantity {
		return fmt.Errorf("quantity must be 1..%d: %w", maxQuantity, core_errors.ErrInvalidArgument)
	}
	return nil
}

type CartItemPatchDTO struct {
	Quantity int `json:"quantity"`
}

func (d *CartItemPatchDTO) Validate() error {
	if d.Quantity < 1 || d.Quantity > maxQuantity {
		return fmt.Errorf("quantity must be 1..%d: %w", maxQuantity, core_errors.ErrInvalidArgument)
	}
	return nil
}

type CartItemResponseDTO struct {
	ID           uuid.UUID `json:"id"`
	DishID       uuid.UUID `json:"dish_id"`
	Name         string    `json:"name"`
	PriceKopecks int64     `json:"price_kopecks"`
	Quantity     int       `json:"quantity"`
	TotalKopecks int64     `json:"total_kopecks"`
	IsAvailable  bool      `json:"is_available"`
}

type CartResponseDTO struct {
	RestID       *uuid.UUID            `json:"restaurant_id"`
	RestName     *string               `json:"restaurant_name"`
	Items        []CartItemResponseDTO `json:"items"`
	TotalKopecks int64                 `json:"total_kopecks"`
}

func cartItemToDomain(item CartItemRequestDTO) domain.CartItemRequest {
	return domain.CartItemRequest{
		DishID:   uuid.MustParse(item.DishID),
		Quantity: item.Quantity,
	}
}

func cartToDTO(cart domain.Cart) CartResponseDTO {
	items := make([]CartItemResponseDTO, len(cart.Items))
	for k, v := range cart.Items {
		items[k] = CartItemResponseDTO{
			ID:           v.ID,
			DishID:       v.DishID,
			Name:         v.Name,
			PriceKopecks: v.PriceKopecks,
			Quantity:     v.Quantity,
			TotalKopecks: v.TotalKopecks(),
			IsAvailable:  v.IsAvailable,
		}
	}

	return CartResponseDTO{
		RestID:       cart.RestID,
		RestName:     cart.RestName,
		Items:        items,
		TotalKopecks: cart.TotalKopecks,
	}
}
