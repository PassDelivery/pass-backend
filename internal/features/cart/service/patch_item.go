package cart_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *CartService) PatchItem(
	ctx context.Context,
	itemID uuid.UUID,
	quantity int,
) (domain.Cart, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("get id from context: %w", err)
	}

	if err := s.cartRepository.PatchItem(ctx, userID, itemID, quantity); err != nil {
		return domain.Cart{}, fmt.Errorf("patch item in service: %w", err)
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
