package menegment_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *MenegmentRepository) GetCategories(
	ctx context.Context,
	restID uuid.UUID,
) ([]domain.CategoryResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, restaurant_id, name, sort_order
	FROM dish_categories
	WHERE restaurant_id = $1
	ORDER BY sort_order, name;
	`

	rows, err := r.pool.Query(ctx, query, restID)
	if err != nil {
		return nil, fmt.Errorf("query categories: %w", err)
	}
	defer rows.Close()

	categories := []domain.CategoryResponse{}

	for rows.Next() {
		var category domain.CategoryResponse

		if err := rows.Scan(
			&category.ID,
			&category.RestID,
			&category.Name,
			&category.SortOrder,
		); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}

		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("categories rows: %w", err)
	}

	return categories, nil
}
