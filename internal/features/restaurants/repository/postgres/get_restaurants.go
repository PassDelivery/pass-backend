package restaurants_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
)

func (r *RestaurantsRepository) GetOpenRestaurants(
	ctx context.Context,
) ([]domain.RestaurantCard, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, name, description, address, phone, is_open, photo_url, created_at
	FROM restaurants
	WHERE is_open = true
	ORDER BY created_at DESC;
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query restaurants: %w", err)
	}
	defer rows.Close()

	rests := []domain.RestaurantCard{}

	for rows.Next() {
		var rest domain.RestaurantCard

		if err := rows.Scan(
			&rest.ID,
			&rest.Name,
			&rest.Description,
			&rest.Address,
			&rest.Phone,
			&rest.IsOpen,
			&rest.PhotoURL,
			&rest.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan restaurant: %w", err)
		}

		rests = append(rests, rest)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("restaurants rows: %w", err)
	}

	return rests, nil
}
