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

func (r *CouriersRepository) SetStatus(
	ctx context.Context,
	userID uuid.UUID,
	status domain.CourierStatus,
) (domain.Courier, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Courier{}, fmt.Errorf("failed to begin: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			return
		}
	}()

	queryCourier := `
	INSERT INTO couriers (user_id, status)
	VALUES ($1, $2)
	ON CONFLICT (user_id) DO UPDATE SET status = EXCLUDED.status
	WHERE couriers.status <> 'busy'
	RETURNING id, user_id, transport_type, status;
	`

	queryDropOffers := `
	DELETE FROM deliveries
	WHERE courier_id = $1 AND status = 'offered';
	`

	var courier domain.Courier

	err = tx.QueryRow(ctx, queryCourier, userID, status).Scan(
		&courier.ID,
		&courier.UserID,
		&courier.TransportType,
		&courier.Status,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Courier{}, fmt.Errorf("courier has active delivery: %w", core_errors.ErrConflict)
	}
	if err != nil {
		return domain.Courier{}, fmt.Errorf("upsert courier: %w", err)
	}

	if status == domain.CourierStatusOffline {
		if _, err := tx.Exec(ctx, queryDropOffers, courier.ID); err != nil {
			return domain.Courier{}, fmt.Errorf("drop offers: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.Courier{}, fmt.Errorf("faild commit to database: %w", err)
	}

	return courier, nil
}
