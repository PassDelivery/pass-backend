package cart_service

import (
	"context"
	"fmt"

	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
)

func (s *CartService) ClearCart(
	ctx context.Context,
) error {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get id from context: %w", err)
	}

	if err := s.cartRepository.ClearCart(ctx, userID); err != nil {
		return fmt.Errorf("clear cart in service: %w", err)
	}

	return nil
}
