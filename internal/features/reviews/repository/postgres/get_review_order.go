package reviews_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *ReviewsRepository) GetReviewOrder(
	ctx context.Context,
	orderID uuid.UUID,
) (domain.ReviewOrder, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT user_id, restaurant_id, status
	FROM orders
	WHERE id = $1;
	`

	var order domain.ReviewOrder

	err := r.pool.QueryRow(ctx, query, orderID).Scan(
		&order.UserID,
		&order.RestID,
		&order.Status,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ReviewOrder{}, fmt.Errorf("order not found: %w", core_errors.ErrNotFound)
	}
	if err != nil {
		return domain.ReviewOrder{}, fmt.Errorf("scan order: %w", err)
	}

	return order, nil
}
