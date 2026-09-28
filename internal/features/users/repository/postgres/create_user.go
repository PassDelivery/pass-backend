package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.UserWithHash,
) (domain.UserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO users (email, password_hash, name, phone, role)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, email, name, phone, role, created_at;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.Email,
		user.Hash,
		user.Name,
		user.Phone,
		user.Role,
	)

	var userModel UserModel

	err := row.Scan(
		&userModel.ID,
		&userModel.Email,
		&userModel.Name,
		&userModel.Phone,
		&userModel.Role,
		&userModel.CreatedAt,
	)

	if err != nil {
		return domain.UserResponse{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUserResponse(
		userModel.ID,
		domain.UninitializedToken,
		userModel.Email,
		userModel.Name,
		userModel.Phone,
		userModel.Role,
		userModel.CreatedAt,
	)

	return userDomain, nil
}
