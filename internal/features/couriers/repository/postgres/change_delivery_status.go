package couriers_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *CouriersRepository) ChangeDeliveryStatus(
	ctx context.Context,
	courierID uuid.UUID,
	userID uuid.UUID,
	deliveryID uuid.UUID,
	from domain.DeliveryStatus,
	to domain.DeliveryStatus,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			return
		}
	}()

	queryDelivery := `
	UPDATE deliveries
	SET status = $1, delivered_at = CASE WHEN $1 = 'delivered' THEN NOW() ELSE delivered_at END
	WHERE id = $2 AND courier_id = $3 AND status = $4
	RETURNING order_id;
	`

	queryOrder := `
	UPDATE orders
	SET status = $1, updated_at = NOW()
	WHERE id = $2;
	`

	queryHistory := `
	INSERT INTO order_status_history (order_id, status, changed_by)
	VALUES ($1, $2, $3);
	`

	queryFree := `
	UPDATE couriers
	SET status = 'free'
	WHERE id = $1;
	`

	var orderID uuid.UUID

	err = tx.QueryRow(ctx, queryDelivery, to, deliveryID, courierID, from).Scan(&orderID)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("delivery status already changed: %w", core_errors.ErrConflict)
	}
	if err != nil {
		return fmt.Errorf("update delivery: %w", err)
	}

	orderStatus := to.OrderStatus()

	if _, err := tx.Exec(ctx, queryOrder, orderStatus, orderID); err != nil {
		return fmt.Errorf("update order: %w", err)
	}

	if _, err := tx.Exec(ctx, queryHistory, orderID, orderStatus, userID); err != nil {
		return fmt.Errorf("insert status history: %w", err)
	}

	if to == domain.DeliveryStatusDelivered {
		if _, err := tx.Exec(ctx, queryFree, courierID); err != nil {
			return fmt.Errorf("set courier free: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("faild commit to database: %w", err)
	}

	return nil
}
