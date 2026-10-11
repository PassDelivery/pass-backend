package users_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
)


func (s *UsersService) GetAddresses(ctx context.Context) ([]domain.AddressesResponse, error) {
	id, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get id from context: %w", err)
	}

	addresses, err := s.usersRepository.GetAddresses(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get addresses: %w", err)
	}

	return addresses, nil
}