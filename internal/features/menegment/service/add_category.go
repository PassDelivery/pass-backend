package menegment_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

func (s *MenegmentService) AddCategory(
	ctx context.Context,
	category domain.CategoryRequest,
	restID uuid.UUID,
) (domain.CategoryResponse, error) {
	catResp, err := s.menegmentRepository.AddCategory(ctx, category, restID)
	if err != nil {
		return domain.CategoryResponse{}, fmt.Errorf("add category in service: %w", err)
	}

	return catResp, nil
}
