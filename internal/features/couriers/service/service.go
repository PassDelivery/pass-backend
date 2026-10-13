package couriers_service

import (
	"context"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

type CouriersService struct {
	couriersRepository CouriersRepository
}

type CouriersRepository interface {
	GetCourier(
		ctx context.Context,
		userID uuid.UUID,
	) (domain.Courier, error)
	SetStatus(
		ctx context.Context,
		userID uuid.UUID,
		status domain.CourierStatus,
	) (domain.Courier, error)
	OfferDeliveries(
		ctx context.Context,
		courierID uuid.UUID,
		limit int,
	) error
	GetDeliveries(
		ctx context.Context,
		courierID uuid.UUID,
	) ([]domain.Delivery, error)
	GetDelivery(
		ctx context.Context,
		courierID uuid.UUID,
		deliveryID uuid.UUID,
	) (domain.Delivery, error)
	AcceptDelivery(
		ctx context.Context,
		courierID uuid.UUID,
		deliveryID uuid.UUID,
	) error
	RefuseDelivery(
		ctx context.Context,
		courierID uuid.UUID,
		deliveryID uuid.UUID,
	) error
	ChangeDeliveryStatus(
		ctx context.Context,
		courierID uuid.UUID,
		userID uuid.UUID,
		deliveryID uuid.UUID,
		from domain.DeliveryStatus,
		to domain.DeliveryStatus,
	) error
}

func NewCouriersService(cr CouriersRepository) *CouriersService {
	return &CouriersService{
		couriersRepository: cr,
	}
}
