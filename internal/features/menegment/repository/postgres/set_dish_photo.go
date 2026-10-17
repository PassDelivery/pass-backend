package menegment_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *MenegmentRepository) SetDishPhoto(
	ctx context.Context,
	dishID uuid.UUID,
	photoURL string,
) (domain.DishResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE dishes
	SET photo_url = $1
	WHERE id = $2
	RETURNING id, restaurant_id, category_id, name, description, price_kopecks, cooking_time_minutes, is_available, photo_url, created_at;
	`

	var dish domain.DishResponse

	err := r.pool.QueryRow(ctx, query, photoURL, dishID).Scan(
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
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.DishResponse{}, fmt.Errorf("dish not found: %w", core_errors.ErrNotFound)
	}
	if err != nil {
		return domain.DishResponse{}, fmt.Errorf("scan dish: %w", err)
	}

	return dish, nil
}
