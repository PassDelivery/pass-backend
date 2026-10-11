package orders_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *OrdersRepository) GetOrderStatus(
	ctx context.Context,
	orderID uuid.UUID,
) (domain.OrderStatus, uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT o.status, r.owner_id
	FROM orders o
	JOIN restaurants r ON r.id = o.restaurant_id
	WHERE o.id = $1;
	`

	var (
		status  domain.OrderStatus
		ownerID uuid.UUID
	)

	err := r.pool.QueryRow(ctx, query, orderID).Scan(&status, &ownerID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", uuid.Nil, fmt.Errorf("order not found: %w", core_errors.ErrNotFound)
	}
	if err != nil {
		return "", uuid.Nil, fmt.Errorf("scan error: %w", err)
	}

	return status, ownerID, nil
}
