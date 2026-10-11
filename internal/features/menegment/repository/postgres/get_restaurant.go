package menegment_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *MenegmentRepository) GetRestaurant(
	ctx context.Context,
	restID uuid.UUID,
) (domain.RestaurantRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT name, description, address, phone, is_open
	FROM restaurants
	WHERE id = $1;
	`

	var rest domain.RestaurantRequest

	if err := r.pool.QueryRow(ctx, query, restID).Scan(
		&rest.Name,
		&rest.Description,
		&rest.Address,
		&rest.Phone,
		&rest.IsOpen,
	); err != nil {
		return domain.RestaurantRequest{}, fmt.Errorf("scan error: %w", err)
	}

	return rest, nil
}
