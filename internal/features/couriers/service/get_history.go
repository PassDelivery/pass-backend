package couriers_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
)

const historyLimit = 100

func (s *CouriersService) GetHistory(
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

	deliveries, err := s.couriersRepository.GetHistory(ctx, courier.ID, historyLimit)
	if err != nil {
		return nil, fmt.Errorf("get history in service: %w", err)
	}

	return deliveries, nil
}
