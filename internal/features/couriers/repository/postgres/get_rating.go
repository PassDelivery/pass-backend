package couriers_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *CouriersRepository) GetRating(
	ctx context.Context,
	courierID uuid.UUID,
) (domain.CourierRating, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT AVG(rv.courier_rating)::float8, COUNT(rv.courier_rating)
	FROM reviews rv
	JOIN deliveries d ON d.order_id = rv.order_id
	WHERE d.courier_id = $1 AND d.status = 'delivered' AND rv.courier_rating IS NOT NULL;
	`

	var rating domain.CourierRating

	if err := r.pool.QueryRow(ctx, query, courierID).Scan(
		&rating.Rating,
		&rating.Count,
	); err != nil {
		return domain.CourierRating{}, fmt.Errorf("scan courier rating: %w", err)
	}

	return rating, nil
}
