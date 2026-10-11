package menegment_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *MenegmentService) PatchRestaurant(
	ctx context.Context,
	rest domain.RestaurantPatch,
	restID uuid.UUID,
) (domain.RestaurantResponse, error) {
	ownerID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.RestaurantResponse{}, fmt.Errorf("get id from context: %w", err)
	}

	restReq, err := s.menegmentRepository.GetRestaurant(ctx, restID)
	if err != nil {
		return domain.RestaurantResponse{}, fmt.Errorf("get restaurant: %w", err)
	}

	restReq.ApplyPatch(rest)

	restResp, err := s.menegmentRepository.PatchRestaurant(ctx, restReq, restID, ownerID)
	if err != nil {
		return domain.RestaurantResponse{}, fmt.Errorf("patch restaurant in service: %w", err)
	}

	return restResp, nil
}
