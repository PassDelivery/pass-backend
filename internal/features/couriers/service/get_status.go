package couriers_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
)

func (s *CouriersService) GetStatus(
	ctx context.Context,
) (domain.Courier, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.Courier{}, fmt.Errorf("get id from context: %w", err)
	}

	courier, err := s.couriersRepository.GetCourier(ctx, userID)
	if errors.Is(err, core_errors.ErrNotFound) {
		return domain.Courier{UserID: userID, Status: domain.CourierStatusOffline}, nil
	}
	if err != nil {
		return domain.Courier{}, fmt.Errorf("get courier in service: %w", err)
	}

	return courier, nil
}
