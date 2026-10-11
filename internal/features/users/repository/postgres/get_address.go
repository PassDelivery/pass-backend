package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *UsersRepository) GetAddress(
	ctx context.Context,
	addrID uuid.UUID,
) (domain.AddressRequest, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT city, street, house, apartment, comment, is_default
	FROM addresses
	WHERE id = $1;
	`

	var address domain.AddressRequest

	if err := r.pool.QueryRow(ctx, query, addrID).Scan(
		&address.City,
		&address.Street,
		&address.House,
		&address.Apartment,
		&address.Comment,
		&address.IsDefault,
	); err != nil {
		return domain.AddressRequest{}, fmt.Errorf("scan error: %w", err)
	}

	return address, nil
}
