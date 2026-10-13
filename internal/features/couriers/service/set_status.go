package couriers_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
)

func (s *CouriersService) SetStatus(
	ctx context.Context,
	status domain.CourierStatus,
) (domain.Courier, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.Courier{}, fmt.Errorf("get id from context: %w", err)
	}

	if !status.CanBeSetByCourier() {
		return domain.Courier{}, fmt.Errorf("courier cant set status %s: %w", status, core_errors.ErrInvalidArgument)
	}

	courier, err := s.couriersRepository.SetStatus(ctx, userID, status)
	if err != nil {
		return domain.Courier{}, fmt.Errorf("set status in service: %w", err)
	}

	return courier, nil
}
