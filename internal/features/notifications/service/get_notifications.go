package notifications_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
)

const notificationsLimit = 50

func (s *NotificationsService) GetNotifications(
	ctx context.Context,
) ([]domain.NotificationResponse, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get id from context: %w", err)
	}

	notifications, err := s.notificationsRepository.GetNotifications(ctx, userID, notificationsLimit)
	if err != nil {
		return nil, fmt.Errorf("get notifications in service: %w", err)
	}

	return notifications, nil
}
