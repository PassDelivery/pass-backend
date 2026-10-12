package cart_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
)

func (s *CartService) AddItem(
	ctx context.Context,
	item domain.CartItemRequest,
) (domain.Cart, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.Cart{}, fmt.Errorf("get id from context: %w", err)
	}

	if err := s.cartRepository.AddItem(ctx, userID, item); err != nil {
		return domain.Cart{}, fmt.Errorf("add item in service: %w", err)
	}

	return s.cartWithTotal(ctx, userID)
}
