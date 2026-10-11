package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)


func (r *UsersRepository) GetAddresses(
	ctx context.Context, 
	id uuid.UUID,
	) ([]domain.AddressesResponse, error) {
		ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
		defer cancel()

		query := `
		SELECT addresses.city, addresses.street, addresses.house, addresses.apartment,
					addresses.comment, addresses.is_default, addresses.id, addresses.user_id,
					addresses.created_at
		FROM addresses
		INNER JOIN users ON users.id = addresses.user_id
		WHERE users.id = $1;
		`


		rows, err := r.pool.Query(ctx, query, id)
		if err != nil {
			return nil, fmt.Errorf("query error: %w", err)
		}
		defer rows.Close()

		var addressesModels []AddressesResponseModel
		for rows.Next() {
			
			var addr AddressesResponseModel

			err := rows.Scan(
				&addr.City,
				&addr.Street,
				&addr.House,
				&addr.Apartment,
				&addr.Comment,
				&addr.IsDefault,
				&addr.ID,
				&addr.UserID,
				&addr.CreatedAt,
			)
			if err != nil {
				return nil, fmt.Errorf("scan rows: %w", err)
			}

			addressesModels = append(addressesModels, addr)
		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("next rows: %w", err)
		}

		addressResponse := addressDomainsFromModels(addressesModels)

		return addressResponse, nil
}