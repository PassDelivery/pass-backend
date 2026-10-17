package menegment_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *MenegmentRepository) GetDishes(
	ctx context.Context,
	restID uuid.UUID,
) ([]domain.DishResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, restaurant_id, category_id, name, description, price_kopecks, cooking_time_minutes, is_available, photo_url, created_at
	FROM dishes
	WHERE restaurant_id = $1
	ORDER BY created_at;
	`

	rows, err := r.pool.Query(ctx, query, restID)
	if err != nil {
		return nil, fmt.Errorf("query dishes: %w", err)
	}
	defer rows.Close()

	dishes := []domain.DishResponse{}

	for rows.Next() {
		var dish domain.DishResponse

		if err := rows.Scan(
			&dish.ID,
			&dish.RestID,
			&dish.CategoryID,
			&dish.Name,
			&dish.Description,
			&dish.PriceKopecks,
			&dish.CookingTimeMinutes,
			&dish.IsAvailable,
			&dish.PhotoURL,
			&dish.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan dish: %w", err)
		}

		dishes = append(dishes, dish)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("dishes rows: %w", err)
	}

	return dishes, nil
}
