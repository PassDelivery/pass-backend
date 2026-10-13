package notifications_service

import (
	"context"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

type NotificationsService struct {
	notificationsRepository NotificationsRepository
}

type NotificationsRepository interface {
	GetNotifications(
		ctx context.Context,
		userID uuid.UUID,
		limit int,
	) ([]domain.NotificationResponse, error)
	ReadNotification(
		ctx context.Context,
		notificationID uuid.UUID,
		userID uuid.UUID,
	) (domain.NotificationResponse, error)
}

func NewNotificationsService(nr NotificationsRepository) *NotificationsService {
	return &NotificationsService{
		notificationsRepository: nr,
	}
}
