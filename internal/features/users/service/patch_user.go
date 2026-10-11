package users_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	"github.com/google/uuid"
)

func (s *UsersService) PatchUser(
	ctx context.Context,
	id uuid.UUID,
	user domain.UserPatch,
) (domain.GetUserResponse, error) {
	log := core_logger.FromContext(ctx)
	log.Debug("invoke PathUser service")
	userRes, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.GetUserResponse{}, fmt.Errorf("get user failed: %w", err)
	}

	if err := userRes.ApplyPatch(user); err != nil {
		return domain.GetUserResponse{}, fmt.Errorf("failed to apply patch: %w", err)
	}

	userRes, err = s.usersRepository.PatchUser(ctx, id, userRes)
	if err != nil {
		return domain.GetUserResponse{}, fmt.Errorf("failed patch user: %w", err)
	}


	return userRes, nil
}
