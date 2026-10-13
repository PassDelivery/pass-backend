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

func (r *CouriersRepository) RefuseDelivery(
	ctx context.Context,
	courierID uuid.UUID,
	deliveryID uuid.UUID,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	queryStatus := `
	SELECT status
	FROM deliveries
	WHERE id = $1 AND courier_id = $2;
	`

	queryRefuse := `
	UPDATE deliveries
	SET status = 'refused', responded_at = NOW()
	WHERE id = $1 AND courier_id = $2 AND status = 'offered';
	`

	var status domain.DeliveryStatus

	err := r.pool.QueryRow(ctx, queryStatus, deliveryID, courierID).Scan(&status)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("delivery not found: %w", core_errors.ErrNotFound)
	}
	if err != nil {
		return fmt.Errorf("scan delivery: %w", err)
	}
	if status != domain.DeliveryStatusOffered {
		return fmt.Errorf("delivery is %s: %w", status, core_errors.ErrConflict)
	}

	tag, err := r.pool.Exec(ctx, queryRefuse, deliveryID, courierID)
	if err != nil {
		return fmt.Errorf("failed exec: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delivery status already changed: %w", core_errors.ErrConflict)
	}

	return nil
}
