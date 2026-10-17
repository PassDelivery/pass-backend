package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (r *UsersRepository) CreateAddress(
	ctx context.Context,
	address domain.AddressRequest,
	id uuid.UUID,
) (domain.AddressesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.AddressesResponse{}, fmt.Errorf("failed to begin: %w", err)
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			return
		}
	} ()

	queryPatch := `
	UPDATE addresses 
	SET is_default = false
	WHERE user_id = $1 AND is_default = true;
	`

	queryPost := `
	INSERT INTO addresses (user_id, city, street, house, apartment, comment, is_default)
	VALUES ($1, $2, $3, $4, $5, $6, $7 OR NOT EXISTS (SELECT 1 FROM addresses WHERE user_id = $1))
	RETURNING id, user_id, city, street, house, apartment, comment, is_default, created_at;
	`

	var response AddressesResponseModel

	if address.IsDefault {
		_, err = tx.Exec(ctx, queryPatch, id)
		if err != nil {
			return domain.AddressesResponse{}, fmt.Errorf("exec to database: %w", err)
		}
	}

	if err := tx.QueryRow(
		ctx,
		queryPost,
		id,
		address.City,
		address.Street,
		address.House,
		address.Apartment,
		address.Comment,
		address.IsDefault,
	).Scan(
		&response.ID,
		&response.UserID,
		&response.City,
		&response.Street,
		&response.House,
		&response.Apartment,
		&response.Comment,
		&response.IsDefault,
		&response.CreatedAt,
	); err != nil {
		return domain.AddressesResponse{}, fmt.Errorf("scan row: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return domain.AddressesResponse{}, fmt.Errorf("faild commit to database: %w", err)
	}

	return addressDomainFromModel(response), nil
}
