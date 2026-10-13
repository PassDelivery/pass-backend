package orders_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *OrdersRepository) GetStatusHistory(
	ctx context.Context,
	orderID uuid.UUID,
) ([]domain.OrderStatusChange, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT status, changed_by, changed_at
	FROM order_status_history
	WHERE order_id = $1
	ORDER BY changed_at;
	`

	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("query status history: %w", err)
	}
	defer rows.Close()

	history := []domain.OrderStatusChange{}

	for rows.Next() {
		var change domain.OrderStatusChange

		if err := rows.Scan(
			&change.Status,
			&change.ChangedBy,
			&change.ChangedAt,
		); err != nil {
			return nil, fmt.Errorf("scan status change: %w", err)
		}

		history = append(history, change)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("status history rows: %w", err)
	}

	return history, nil
}
