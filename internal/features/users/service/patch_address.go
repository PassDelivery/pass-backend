package users_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *UsersService) PatchAddress(
	ctx context.Context,
	address domain.AddressPatch,
	addrID uuid.UUID,
) (domain.AddressesResponse, error) {

	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.AddressesResponse{}, fmt.Errorf("get id from context: %w", err)
	}

	addr, err := s.usersRepository.GetAddress(ctx, addrID)
	if err != nil {
		return domain.AddressesResponse{}, fmt.Errorf("get address: %w", err)
	}

	if err := addr.ApplyPatch(address); err != nil {
		return domain.AddressesResponse{}, fmt.Errorf("apply path to address: %w", err)
	}
	addressResp, err := s.usersRepository.PatchAddress(ctx, addr, userID, addrID)

	if err != nil {
		return domain.AddressesResponse{}, fmt.Errorf("failed to patch address in repository: %w", err)
	}

	return addressResp, nil
}
