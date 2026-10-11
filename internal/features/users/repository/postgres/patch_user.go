package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	"github.com/google/uuid"
)

func (r *UsersRepository) PatchUser(
	ctx context.Context, id uuid.UUID,
	user domain.GetUserResponse,
) (domain.GetUserResponse, error) {
	log := core_logger.FromContext(ctx)
	log.Debug("invoke PatchUser repository")
	
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	UPDATE users
	SET email=$1, name=$2, phone=$3
	WHERE id=$4
	RETURNING id, email, name, phone, role, created_at;
	`

	var response domain.GetUserResponse
	if err := r.pool.QueryRow(
		ctx, query, 
		user.Email, 
		user.Name, 
		user.Phone, 
		user.ID,
		).Scan(
			&response.ID,
			&response.Email,
			&response.Name,
			&response.Phone,
			&response.Role,
			&response.CreatedAt,
		); err != nil {
			return domain.GetUserResponse{}, fmt.Errorf("failed to scan: %w", err)
	}

	return response, nil
}
