package menegment_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *MenegmentRepository) AddDish(
	ctx context.Context,
	dish domain.DishRequest,
	restID uuid.UUID,
) (domain.DishResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO dishes (restaurant_id, category_id, name, description, price_kopecks, cooking_time_minutes, is_available)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, restaurant_id, category_id, name, description, price_kopecks, cooking_time_minutes, is_available, created_at;
	`

	var dishResp domain.DishResponse

	if err := r.pool.QueryRow(
		ctx,
		query,
		restID,
		dish.CategoryID,
		dish.Name,
		dish.Description,
		dish.PriceKopecks,
		dish.CookingTimeMinutes,
		dish.IsAvailable,
	).Scan(
		&dishResp.ID,
		&dishResp.RestID,
		&dishResp.CategoryID,
		&dishResp.Name,
		&dishResp.Description,
		&dishResp.PriceKopecks,
		&dishResp.CookingTimeMinutes,
		&dishResp.IsAvailable,
		&dishResp.CreatedAt,
	); err != nil {
		return domain.DishResponse{}, fmt.Errorf("failed scan: %w", err)
	}

	return dishResp, nil
}
