package menegment_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *MenegmentRepository) GetRestaurants(
	ctx context.Context,
	ownerID uuid.UUID,
) ([]domain.RestaurantResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT name, description, address, phone, is_open, photo_url, id, owner_id, created_at
	FROM restaurants
	WHERE owner_id = $1
	ORDER BY created_at;
	`

	rows, err := r.pool.Query(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("query restaurants: %w", err)
	}
	defer rows.Close()

	restaurants := []domain.RestaurantResponse{}

	for rows.Next() {
		var restaurant domain.RestaurantResponse

		if err := rows.Scan(
			&restaurant.Name,
			&restaurant.Description,
			&restaurant.Address,
			&restaurant.Phone,
			&restaurant.IsOpen,
			&restaurant.PhotoURL,
			&restaurant.ID,
			&restaurant.OwnerID,
			&restaurant.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan restaurant: %w", err)
		}

		restaurants = append(restaurants, restaurant)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("restaurants rows: %w", err)
	}

	return restaurants, nil
}
