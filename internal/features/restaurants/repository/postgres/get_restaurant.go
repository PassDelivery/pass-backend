package restaurants_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *RestaurantsRepository) GetRestaurant(
	ctx context.Context,
	restID uuid.UUID,
) (domain.RestaurantCard, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, name, description, address, phone, is_open, created_at
	FROM restaurants
	WHERE id = $1;
	`

	var rest domain.RestaurantCard

	err := r.pool.QueryRow(ctx, query, restID).Scan(
		&rest.ID,
		&rest.Name,
		&rest.Description,
		&rest.Address,
		&rest.Phone,
		&rest.IsOpen,
		&rest.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.RestaurantCard{}, fmt.Errorf("restaurant not found: %w", core_errors.ErrNotFound)
	}
	if err != nil {
		return domain.RestaurantCard{}, fmt.Errorf("scan error: %w", err)
	}

	return rest, nil
}
