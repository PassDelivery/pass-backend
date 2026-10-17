package menegment_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *MenegmentRepository) PatchRestaurant(
	ctx context.Context,
	rest domain.RestaurantRequest,
	restID uuid.UUID,
	ownerID uuid.UUID,
) (domain.RestaurantResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE restaurants
	SET name=$1, description=$2, address=$3, phone=$4, is_open=$5
	WHERE owner_id=$6 AND id=$7
	RETURNING name, description, address, phone, is_open, photo_url, id, owner_id, created_at;
	`

	var restResp domain.RestaurantResponse

	if err := r.pool.QueryRow(
		ctx,
		query,
		rest.Name,
		rest.Description,
		rest.Address,
		rest.Phone,
		rest.IsOpen,
		ownerID,
		restID,
	).Scan(
		&restResp.Name,
		&restResp.Description,
		&restResp.Address,
		&restResp.Phone,
		&restResp.IsOpen,
		&restResp.PhotoURL,
		&restResp.ID,
		&restResp.OwnerID,
		&restResp.CreatedAt,
	); err != nil {
		return domain.RestaurantResponse{}, fmt.Errorf("failed scan: %w", err)
	}

	return restResp, nil
}
