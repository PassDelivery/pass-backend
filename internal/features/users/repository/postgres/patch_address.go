package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *UsersRepository) PatchAddress(
	ctx context.Context,
	address domain.AddressRequest,
	userID uuid.UUID,
	addrID uuid.UUID,
) (domain.AddressesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE addresses
	SET city=$1, street=$2, house=$3, apartment=$4, comment=$5, is_default=$6
	WHERE user_id=$7 AND id=$8
	RETURNING id, user_id, city, street, house, apartment, comment, is_default, created_at;
	`

	var addressModel AddressesResponseModel

	if err := r.pool.QueryRow(
		ctx, 
		query,
		address.City,
		address.Street,
		address.House,
		address.Apartment,
		address.Comment,
		address.IsDefault,
		userID,
		addrID,
	).Scan(
		&addressModel.ID,
		&addressModel.UserID,
		&addressModel.City,
		&addressModel.Street,
		&addressModel.House,
		&addressModel.Apartment,
		&addressModel.Comment,
		&addressModel.IsDefault,
		&addressModel.CreatedAt,
	); err != nil {
		return domain.AddressesResponse{}, fmt.Errorf("failed to scan: %w", err)
	}

	return addressDomainFromModel(addressModel), nil
}
