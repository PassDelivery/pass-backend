package couriers_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *CouriersRepository) OfferDeliveries(
	ctx context.Context,
	courierID uuid.UUID,
	limit int,
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

	queryLock := `
	SELECT status
	FROM couriers
	WHERE id = $1
	FOR UPDATE;
	`

	queryDropStale := `
	DELETE FROM deliveries d
	USING orders o
	WHERE d.order_id = o.id AND d.courier_id = $1 AND d.status = 'offered' AND o.status NOT IN ('cooking', 'ready');
	`

	queryOffer := `
	INSERT INTO deliveries (order_id, courier_id, status)
	SELECT o.id, $1, 'offered'
	FROM orders o
	WHERE o.status IN ('cooking', 'ready')
	AND NOT EXISTS (
		SELECT 1 FROM deliveries d
		WHERE d.order_id = o.id
		AND (d.status IN ('accepted', 'picked_up', 'delivered') OR d.courier_id = $1)
	)
	ORDER BY o.updated_at
	LIMIT $2;
	`

	var status domain.CourierStatus

	if err := tx.QueryRow(ctx, queryLock, courierID).Scan(&status); err != nil {
		return fmt.Errorf("lock courier: %w", err)
	}

	if _, err := tx.Exec(ctx, queryDropStale, courierID); err != nil {
		return fmt.Errorf("drop stale offers: %w", err)
	}

	if status == domain.CourierStatusFree {
		if _, err := tx.Exec(ctx, queryOffer, courierID, limit); err != nil {
			return fmt.Errorf("offer deliveries: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("faild commit to database: %w", err)
	}

	return nil
}
