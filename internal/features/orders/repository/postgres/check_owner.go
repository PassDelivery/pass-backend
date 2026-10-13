package orders_postgres_repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *OrdersRepository) IsRestaurantOwner(
	ctx context.Context,
	restID uuid.UUID,
	ownerID uuid.UUID,
) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT EXISTS(
		SELECT 1 FROM restaurants
		WHERE id = $1 AND owner_id = $2
	);
	`

	var ok bool

	if err := r.pool.QueryRow(ctx, query, restID, ownerID).Scan(&ok); err != nil {
		return false, fmt.Errorf("failed scan: %w", err)
	}

	return ok, nil
}
