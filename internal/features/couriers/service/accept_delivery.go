package couriers_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *CouriersService) AcceptDelivery(
	ctx context.Context,
	deliveryID uuid.UUID,
) (domain.Delivery, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.Delivery{}, fmt.Errorf("get id from context: %w", err)
	}

	courier, err := s.couriersRepository.GetCourier(ctx, userID)
	if err != nil {
		return domain.Delivery{}, fmt.Errorf("get courier: %w", err)
	}

	if err := s.couriersRepository.AcceptDelivery(ctx, courier.ID, deliveryID); err != nil {
		return domain.Delivery{}, fmt.Errorf("accept delivery in service: %w", err)
	}

	delivery, err := s.couriersRepository.GetDelivery(ctx, courier.ID, deliveryID)
	if err != nil {
		return domain.Delivery{}, fmt.Errorf("get delivery: %w", err)
	}

	return delivery, nil
}
