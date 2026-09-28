package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)


func (r *UsersRepository) GetUser(ctx context.Context, id uuid.UUID) (domain.GetUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, email, name, phone, role, created_at
	FROM users
	WHERE id = $1;
	`

	var user domain.GetUserResponse

	if err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
	); err != nil {
		return domain.GetUserResponse{}, fmt.Errorf("scan error: %w", err)
	}


	return user, nil
}