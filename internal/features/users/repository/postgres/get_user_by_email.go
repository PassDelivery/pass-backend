package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
)

func (r *UsersRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (domain.UserResponse, string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, email, password_hash, name, phone, role, created_at
	FROM users
	WHERE LOWER(email) = LOWER($1);
	`

	var user UserModelWithHash

	if err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Hash,
		&user.Name,
		&user.Phone,
		&user.Role,
		&user.CreatedAt,
	); err != nil {
		return domain.UserResponse{}, "", fmt.Errorf("failed scan from database: %w", err)
	}
	hash := user.Hash

	userRes := domain.UserResponse{
		ID:        user.ID,
		Token:     domain.UninitializedToken,
		Email:     user.Email,
		Name:      user.Name,
		Phone:     user.Phone,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	return userRes, hash, nil
}
