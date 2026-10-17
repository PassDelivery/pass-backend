package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) SetAvatar(
	ctx context.Context,
	userID uuid.UUID,
	avatarURL string,
) (domain.GetUserResponse, *string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	WITH old AS (
		SELECT avatar_url FROM users WHERE id = $2
	)
	UPDATE users
	SET avatar_url = $1
	WHERE id = $2
	RETURNING id, email, name, phone, role, avatar_url, created_at, (SELECT avatar_url FROM old);
	`

	var (
		user   domain.GetUserResponse
		oldURL *string
	)

	err := r.pool.QueryRow(ctx, query, avatarURL, userID).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Phone,
		&user.Role,
		&user.AvatarURL,
		&user.CreatedAt,
		&oldURL,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.GetUserResponse{}, nil, fmt.Errorf("user not found: %w", core_errors.ErrNotFound)
	}
	if err != nil {
		return domain.GetUserResponse{}, nil, fmt.Errorf("scan user: %w", err)
	}

	return user, oldURL, nil
}
