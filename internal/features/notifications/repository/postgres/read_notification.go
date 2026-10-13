package notifications_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *NotificationsRepository) ReadNotification(
	ctx context.Context,
	notificationID uuid.UUID,
	userID uuid.UUID,
) (domain.NotificationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE notifications
	SET is_read = true
	WHERE id = $1 AND user_id = $2
	RETURNING id, order_id, type, body, is_read, created_at;
	`

	var notification domain.NotificationResponse

	err := r.pool.QueryRow(ctx, query, notificationID, userID).Scan(
		&notification.ID,
		&notification.OrderID,
		&notification.Type,
		&notification.Body,
		&notification.IsRead,
		&notification.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.NotificationResponse{}, fmt.Errorf("notification not found: %w", core_errors.ErrNotFound)
	}
	if err != nil {
		return domain.NotificationResponse{}, fmt.Errorf("failed scan: %w", err)
	}

	return notification, nil
}
