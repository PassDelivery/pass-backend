package menegment_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *MenegmentRepository) AddCategory(
	ctx context.Context,
	category domain.CategoryRequest,
	restID uuid.UUID,
) (domain.CategoryResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO dish_categories (name, sort_order, restaurant_id)
	VALUES ($1, $2, $3)
	RETURNING id, restaurant_id, name, sort_order;
	`

	var catResp domain.CategoryResponse

	if err := r.pool.QueryRow(
		ctx,
		query,
		category.Name,
		category.SortOrder,
		restID,
	).Scan(
		&catResp.ID,
		&catResp.RestID,
		&catResp.Name,
		&catResp.SortOrder,
	); err != nil {
		return domain.CategoryResponse{}, fmt.Errorf("failed scan: %w", err)
	}

	return catResp, nil
}
