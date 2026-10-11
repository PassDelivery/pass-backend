package users_service

import (
	"context"
	"fmt"

	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *UsersService) DeleteAddress(
	ctx context.Context,
	addrID uuid.UUID,
) error {
	user_id, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get id from context: %w", err)
	}
	if err := s.usersRepository.DeleteAddress(ctx, addrID, user_id); err != nil {
		return fmt.Errorf("delete address in service: %w", err)
	}

	return nil
}
