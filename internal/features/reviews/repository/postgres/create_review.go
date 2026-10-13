package reviews_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *ReviewsRepository) CreateReview(
	ctx context.Context,
	review domain.ReviewRequest,
	orderID uuid.UUID,
	userID uuid.UUID,
	restID uuid.UUID,
) (domain.ReviewResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	WITH rv AS (
		INSERT INTO reviews (order_id, user_id, restaurant_id, rating, courier_rating, body)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, order_id, user_id, restaurant_id, rating, courier_rating, body, created_at
	)
	SELECT rv.id, rv.order_id, rv.restaurant_id, split_part(u.name, ' ', 1), rv.rating, rv.courier_rating, rv.body, rv.created_at
	FROM rv
	JOIN users u ON u.id = rv.user_id;
	`

	var reviewResp domain.ReviewResponse

	err := r.pool.QueryRow(
		ctx,
		query,
		orderID,
		userID,
		restID,
		review.Rating,
		review.CourierRating,
		review.Body,
	).Scan(
		&reviewResp.ID,
		&reviewResp.OrderID,
		&reviewResp.RestID,
		&reviewResp.AuthorName,
		&reviewResp.Rating,
		&reviewResp.CourierRating,
		&reviewResp.Body,
		&reviewResp.CreatedAt,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ReviewResponse{}, fmt.Errorf("review already exists: %w", core_errors.ErrConflict)
	}
	if err != nil {
		return domain.ReviewResponse{}, fmt.Errorf("failed scan: %w", err)
	}

	return reviewResp, nil
}
