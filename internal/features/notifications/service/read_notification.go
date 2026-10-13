package notifications_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *NotificationsService) ReadNotification(
	ctx context.Context,
	notificationID uuid.UUID,
) (domain.NotificationResponse, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.NotificationResponse{}, fmt.Errorf("get id from context: %w", err)
	}

	notification, err := s.notificationsRepository.ReadNotification(ctx, notificationID, userID)
	if err != nil {
		return domain.NotificationResponse{}, fmt.Errorf("read notification in service: %w", err)
	}

	return notification, nil
}
