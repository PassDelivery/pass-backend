package users_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
)

func (r *UsersService) GetUser(ctx context.Context) (domain.GetUserResponse, error) {
	id, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.GetUserResponse{}, fmt.Errorf("get id from context: %w", err)
	}
	user, err := r.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.GetUserResponse{}, fmt.Errorf("get user from pepository: %w", err)
	}

	return user, nil
}
