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

func (r *MenegmentRepository) SetRestaurantPhoto(
	ctx context.Context,
	restID uuid.UUID,
	photoURL string,
) (domain.RestaurantResponse, *string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	WITH old AS (
		SELECT photo_url FROM restaurants WHERE id = $2
	)
	UPDATE restaurants
	SET photo_url = $1
	WHERE id = $2
	RETURNING name, description, address, phone, is_open, photo_url, id, owner_id, created_at, (SELECT photo_url FROM old);
	`

	var (
		restaurant domain.RestaurantResponse
		oldURL     *string
	)

	err := r.pool.QueryRow(ctx, query, photoURL, restID).Scan(
		&restaurant.Name,
		&restaurant.Description,
		&restaurant.Address,
		&restaurant.Phone,
		&restaurant.IsOpen,
		&restaurant.PhotoURL,
		&restaurant.ID,
		&restaurant.OwnerID,
		&restaurant.CreatedAt,
		&oldURL,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.RestaurantResponse{}, nil, fmt.Errorf("restaurant not found: %w", core_errors.ErrNotFound)
	}
	if err != nil {
		return domain.RestaurantResponse{}, nil, fmt.Errorf("scan restaurant: %w", err)
	}

	return restaurant, oldURL, nil
}
