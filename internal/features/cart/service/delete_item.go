package cart_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *CartService) DeleteItem(
	ctx context.Context,
	itemID uuid.UUID,
) (domain.Cart, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("get id from context: %w", err)
	}

	if err := s.cartRepository.DeleteItem(ctx, userID, itemID); err != nil {
		return domain.Cart{}, fmt.Errorf("delete item in service: %w", err)
	}

	cart, err := s.cartRepository.GetCart(ctx, userID)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("get cart in service: %w", err)
	}

	for _, item := range cart.Items {
		cart.TotalKopecks += item.TotalKopecks()
	}

	return cart, nil
}
