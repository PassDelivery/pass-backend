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

func (r *CouriersRepository) GetCourier(
	ctx context.Context,
	userID uuid.UUID,
) (domain.Courier, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, user_id, transport_type, status
	FROM couriers
	WHERE user_id = $1;
	`

	var courier domain.Courier

	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&courier.ID,
		&courier.UserID,
		&courier.TransportType,
		&courier.Status,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Courier{}, fmt.Errorf("courier not found: %w", core_errors.ErrNotFound)
	}
	if err != nil {
		return domain.Courier{}, fmt.Errorf("scan courier: %w", err)
	}

	return courier, nil
}
