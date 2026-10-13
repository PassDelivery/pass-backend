package notifications_transport_http

import (
	"time"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

type NotificationResponseDTO struct {
	ID        uuid.UUID  `json:"id"`
	OrderID   *uuid.UUID `json:"order_id"`
	Type      string     `json:"type"`
	Body      string     `json:"body"`
	IsRead    bool       `json:"is_read"`
	CreatedAt time.Time  `json:"created_at"`
}

func notificationToDTO(notification domain.NotificationResponse) NotificationResponseDTO {
	return NotificationResponseDTO{
		ID:        notification.ID,
		OrderID:   notification.OrderID,
		Type:      notification.Type,
		Body:      notification.Body,
		IsRead:    notification.IsRead,
		CreatedAt: notification.CreatedAt,
	}
}

func notificationsToDTO(notifications []domain.NotificationResponse) []NotificationResponseDTO {
	res := make([]NotificationResponseDTO, len(notifications))
	for k, v := range notifications {
		res[k] = notificationToDTO(v)
	}

	return res
}
