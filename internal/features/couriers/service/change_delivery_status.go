package couriers_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *CouriersService) ChangeDeliveryStatus(
	ctx context.Context,
	deliveryID uuid.UUID,
	status domain.DeliveryStatus,
) (domain.Delivery, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.Delivery{}, fmt.Errorf("get id from context: %w", err)
	}

	courier, err := s.couriersRepository.GetCourier(ctx, userID)
	if err != nil {
		return domain.Delivery{}, fmt.Errorf("get courier: %w", err)
	}

	delivery, err := s.couriersRepository.GetDelivery(ctx, courier.ID, deliveryID)
	if err != nil {
		return domain.Delivery{}, fmt.Errorf("get delivery: %w", err)
	}

	if !delivery.Status.CanCourierMoveTo(status) {
		return domain.Delivery{}, fmt.Errorf("cant move delivery from %s to %s: %w", delivery.Status, status, core_errors.ErrConflict)
	}

	if status == domain.DeliveryStatusPickedUp && delivery.OrderStatus != domain.OrderStatusReady {
		return domain.Delivery{}, fmt.Errorf("order is still %s: %w", delivery.OrderStatus, core_errors.ErrConflict)
	}

	if err := s.couriersRepository.ChangeDeliveryStatus(ctx, courier.ID, userID, deliveryID, delivery.Status, status); err != nil {
		return domain.Delivery{}, fmt.Errorf("change delivery status in service: %w", err)
	}

	delivery, err = s.couriersRepository.GetDelivery(ctx, courier.ID, deliveryID)
	if err != nil {
		return domain.Delivery{}, fmt.Errorf("get delivery: %w", err)
	}

	return delivery, nil
}
