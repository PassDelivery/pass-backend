package restaurants_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *RestaurantsRepository) GetCategories(
	ctx context.Context,
	restID uuid.UUID,
) ([]domain.MenuCategory, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, name, sort_order
	FROM dish_categories
	WHERE restaurant_id = $1
	ORDER BY sort_order NULLS LAST, name;
	`

	rows, err := r.pool.Query(ctx, query, restID)
	if err != nil {
		return nil, fmt.Errorf("query categories: %w", err)
	}
	defer rows.Close()

	cats := []domain.MenuCategory{}

	for rows.Next() {
		var cat domain.MenuCategory

		if err := rows.Scan(
			&cat.ID,
			&cat.Name,
			&cat.SortOrder,
		); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}

		cats = append(cats, cat)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("categories rows: %w", err)
	}

	return cats, nil
}

func (r *RestaurantsRepository) GetAvailableDishes(
	ctx context.Context,
	restID uuid.UUID,
) ([]domain.MenuDish, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, category_id, name, description, price_kopecks, cooking_time_minutes, photo_url
	FROM dishes
	WHERE restaurant_id = $1 AND is_available = true
	ORDER BY created_at;
	`

	rows, err := r.pool.Query(ctx, query, restID)
	if err != nil {
		return nil, fmt.Errorf("query dishes: %w", err)
	}
	defer rows.Close()

	dishes := []domain.MenuDish{}

	for rows.Next() {
		var dish domain.MenuDish

		if err := rows.Scan(
			&dish.ID,
			&dish.CategoryID,
			&dish.Name,
			&dish.Description,
			&dish.PriceKopecks,
			&dish.CookingTimeMinutes,
			&dish.PhotoURL,
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
