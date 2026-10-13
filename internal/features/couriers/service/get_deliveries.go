package couriers_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
)

const offersLimit = 5

func (s *CouriersService) GetDeliveries(
	ctx context.Context,
) ([]domain.Delivery, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get id from context: %w", err)
	}

	courier, err := s.couriersRepository.GetCourier(ctx, userID)
	if errors.Is(err, core_errors.ErrNotFound) {
		return []domain.Delivery{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get courier: %w", err)
	}

	if err := s.couriersRepository.OfferDeliveries(ctx, courier.ID, offersLimit); err != nil {
		return nil, fmt.Errorf("offer deliveries: %w", err)
	}

	deliveries, err := s.couriersRepository.GetDeliveries(ctx, courier.ID)
	if err != nil {
		return nil, fmt.Errorf("get deliveries in service: %w", err)
	}

	for k := range deliveries {
		if deliveries[k].Status == domain.DeliveryStatusOffered {
			deliveries[k].CustomerName = ""
			deliveries[k].CustomerPhone = nil
		}
	}

	return deliveries, nil
}
