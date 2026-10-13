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

func (r *CouriersRepository) AcceptDelivery(
	ctx context.Context,
	courierID uuid.UUID,
	deliveryID uuid.UUID,
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

	queryCourier := `
	SELECT status
	FROM couriers
	WHERE id = $1
	FOR UPDATE;
	`

	queryDelivery := `
	SELECT order_id, status
	FROM deliveries
	WHERE id = $1 AND courier_id = $2
	FOR UPDATE;
	`

	queryOrder := `
	SELECT status
	FROM orders
	WHERE id = $1
	FOR UPDATE;
	`

	queryTaken := `
	SELECT EXISTS(
		SELECT 1 FROM deliveries
		WHERE order_id = $1 AND status IN ('accepted', 'picked_up', 'delivered')
	);
	`

	queryAccept := `
	UPDATE deliveries
	SET status = 'accepted', responded_at = NOW()
	WHERE id = $1;
	`

	queryDropOthers := `
	DELETE FROM deliveries
	WHERE status = 'offered' AND id <> $1 AND (order_id = $2 OR courier_id = $3);
	`

	queryBusy := `
	UPDATE couriers
	SET status = 'busy'
	WHERE id = $1;
	`

	var courierStatus domain.CourierStatus

	if err := tx.QueryRow(ctx, queryCourier, courierID).Scan(&courierStatus); err != nil {
		return fmt.Errorf("lock courier: %w", err)
	}
	if courierStatus != domain.CourierStatusFree {
		return fmt.Errorf("courier is %s: %w", courierStatus, core_errors.ErrConflict)
	}

	var (
		orderID        uuid.UUID
		deliveryStatus domain.DeliveryStatus
	)

	err = tx.QueryRow(ctx, queryDelivery, deliveryID, courierID).Scan(&orderID, &deliveryStatus)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("delivery not found: %w", core_errors.ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("lock delivery: %w", err)
	}
	if deliveryStatus != domain.DeliveryStatusOffered {
		return fmt.Errorf("delivery is %s: %w", deliveryStatus, core_errors.ErrConflict)
	}

	var orderStatus domain.OrderStatus

	if err := tx.QueryRow(ctx, queryOrder, orderID).Scan(&orderStatus); err != nil {
		return fmt.Errorf("lock order: %w", err)
	}
	if orderStatus != domain.OrderStatusCooking && orderStatus != domain.OrderStatusReady {
		return fmt.Errorf("order is %s: %w", orderStatus, core_errors.ErrConflict)
	}

	var taken bool

	if err := tx.QueryRow(ctx, queryTaken, orderID).Scan(&taken); err != nil {
		return fmt.Errorf("check taken: %w", err)
	}
	if taken {
		return fmt.Errorf("order already taken by another courier: %w", core_errors.ErrConflict)
	}

	if _, err := tx.Exec(ctx, queryAccept, deliveryID); err != nil {
		return fmt.Errorf("accept delivery: %w", err)
	}

	if _, err := tx.Exec(ctx, queryDropOthers, deliveryID, orderID, courierID); err != nil {
		return fmt.Errorf("drop other offers: %w", err)
	}

	if _, err := tx.Exec(ctx, queryBusy, courierID); err != nil {
		return fmt.Errorf("set courier busy: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("faild commit to database: %w", err)
	}

	return nil
}
