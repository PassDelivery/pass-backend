package reviews_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *ReviewsRepository) GetRestaurantReviews(
	ctx context.Context,
	restID uuid.UUID,
) ([]domain.ReviewResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT rv.id, rv.order_id, rv.restaurant_id, split_part(u.name, ' ', 1), rv.rating, rv.body, rv.created_at
	FROM reviews rv
	JOIN users u ON u.id = rv.user_id
	WHERE rv.restaurant_id = $1
	ORDER BY rv.created_at DESC;
	`

	rows, err := r.pool.Query(ctx, query, restID)
	if err != nil {
		return nil, fmt.Errorf("query reviews: %w", err)
	}
	defer rows.Close()

	reviews := []domain.ReviewResponse{}

	for rows.Next() {
		var review domain.ReviewResponse

		if err := rows.Scan(
			&review.ID,
			&review.OrderID,
			&review.RestID,
			&review.AuthorName,
			&review.Rating,
			&review.Body,
			&review.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan review: %w", err)
		}

		reviews = append(reviews, review)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("reviews rows: %w", err)
	}

	return reviews, nil
}
