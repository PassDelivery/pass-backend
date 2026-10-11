package menegment_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *MenegmentRepository) RegisterRestaurant(
	ctx context.Context,
	restaurant domain.RestaurantRequest,
	id uuid.UUID,
) (domain.RestaurantResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO restaurants (name, description, address, phone, is_open, owner_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING name, description, address, phone, is_open, id, owner_id, created_at;
	`

	var restResp domain.RestaurantResponse

	if err := r.pool.QueryRow(
		ctx, 
		query, 
		restaurant.Name, 
		restaurant.Description, 
		restaurant.Address,
		restaurant.Phone,
		restaurant.IsOpen,
		id,
	).Scan(
		&restResp.Name,
		&restResp.Description,
		&restResp.Address,
		&restResp.Phone,
		&restResp.IsOpen,
		&restResp.ID,
		&restResp.OwnerID,
		&restResp.CreatedAt,
	); err != nil {
		return domain.RestaurantResponse{}, fmt.Errorf("failed scan: %w", err)
	}

	return restResp, nil
}
