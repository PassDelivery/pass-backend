package menegment_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *MenegmentRepository) PatchDish(
	ctx context.Context,
	dish domain.DishRequest,
	dishID uuid.UUID,
) (domain.DishResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE dishes
	SET category_id=$1, name=$2, description=$3, price_kopecks=$4, cooking_time_minutes=$5, is_available=$6
	WHERE id=$7
	RETURNING id, restaurant_id, category_id, name, description, price_kopecks, cooking_time_minutes, is_available, created_at;
	`

	var dishResp domain.DishResponse

	if err := r.pool.QueryRow(
		ctx,
		query,
		dish.CategoryID,
		dish.Name,
		dish.Description,
		dish.PriceKopecks,
		dish.CookingTimeMinutes,
		dish.IsAvailable,
		dishID,
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
