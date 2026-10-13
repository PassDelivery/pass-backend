package couriers_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
)

func (s *CouriersService) GetRating(
	ctx context.Context,
) (domain.CourierRating, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.CourierRating{}, fmt.Errorf("get id from context: %w", err)
	}

	courier, err := s.couriersRepository.GetCourier(ctx, userID)
	if errors.Is(err, core_errors.ErrNotFound) {
		return domain.CourierRating{}, nil
	}
	if err != nil {
		return domain.CourierRating{}, fmt.Errorf("get courier: %w", err)
	}

	rating, err := s.couriersRepository.GetRating(ctx, courier.ID)
	if err != nil {
		return domain.CourierRating{}, fmt.Errorf("get rating in service: %w", err)
	}

	return rating, nil
}
