package users_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
)

func (s *UsersService) CreateAddress(
	ctx context.Context,
	request domain.AddressRequest,
) (domain.AddressesResponse, error) {
	id, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.AddressesResponse{}, fmt.Errorf("failed to get id from context: %w", err)
	}

	address, err := s.usersRepository.CreateAddress(ctx, request, id)
	if err != nil {
		return domain.AddressesResponse{}, fmt.Errorf("failed to create address: %w", err)
	}

	return address, nil
}
