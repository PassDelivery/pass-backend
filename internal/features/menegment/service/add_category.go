package menegment_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *MenegmentService) AddCategory(
	ctx context.Context,
	category domain.CategoryRequest,
	restID uuid.UUID,
) (domain.CategoryResponse, error) {
	ownerID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.CategoryResponse{}, fmt.Errorf("get id from context: %w", err)
	}

	isOwner, err := s.menegmentRepository.IsRestaurantOwner(ctx, restID, ownerID)
	if err != nil {
		return domain.CategoryResponse{}, fmt.Errorf("check restaurant owner: %w", err)
	}
	if !isOwner {
		return domain.CategoryResponse{}, fmt.Errorf("restaurant not found: %w", core_errors.ErrNotFound)
	}

	catResp, err := s.menegmentRepository.AddCategory(ctx, category, restID)
	if err != nil {
		return domain.CategoryResponse{}, fmt.Errorf("add category in service: %w", err)
	}

	return catResp, nil
}
