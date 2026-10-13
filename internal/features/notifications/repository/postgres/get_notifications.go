package notifications_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *NotificationsRepository) GetNotifications(
	ctx context.Context,
	userID uuid.UUID,
	limit int,
) ([]domain.NotificationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, order_id, type, body, is_read, created_at
	FROM notifications
	WHERE user_id = $1
	ORDER BY created_at DESC
	LIMIT $2;
	`

	rows, err := r.pool.Query(ctx, query, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("query notifications: %w", err)
	}
	defer rows.Close()

	notifications := []domain.NotificationResponse{}

	for rows.Next() {
		var notification domain.NotificationResponse

		if err := rows.Scan(
			&notification.ID,
			&notification.OrderID,
			&notification.Type,
			&notification.Body,
			&notification.IsRead,
			&notification.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}

		notifications = append(notifications, notification)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("notifications rows: %w", err)
	}

	return notifications, nil
}
