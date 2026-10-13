package couriers_service

import (
	"context"
	"fmt"

	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *CouriersService) RefuseDelivery(
	ctx context.Context,
	deliveryID uuid.UUID,
) error {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get id from context: %w", err)
	}

	courier, err := s.couriersRepository.GetCourier(ctx, userID)
	if err != nil {
		return fmt.Errorf("get courier: %w", err)
	}

	if err := s.couriersRepository.RefuseDelivery(ctx, courier.ID, deliveryID); err != nil {
		return fmt.Errorf("refuse delivery in service: %w", err)
	}

	return nil
}
