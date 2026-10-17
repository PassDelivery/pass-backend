package menegment_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_files "github.com/PassDelivery/pass-backend/internal/core/files"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *MenegmentService) UploadRestaurantPhoto(
	ctx context.Context,
	restID uuid.UUID,
	data []byte,
) (domain.RestaurantResponse, error) {
	ownerID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.RestaurantResponse{}, fmt.Errorf("get id from context: %w", err)
	}

	isOwner, err := s.menegmentRepository.IsRestaurantOwner(ctx, restID, ownerID)
	if err != nil {
		return domain.RestaurantResponse{}, fmt.Errorf("check restaurant owner: %w", err)
	}
	if !isOwner {
		return domain.RestaurantResponse{}, fmt.Errorf("restaurant not found: %w", core_errors.ErrNotFound)
	}

	url, err := core_files.SavePhoto(core_files.PhotoKindRestaurants, data)
	if err != nil {
		return domain.RestaurantResponse{}, fmt.Errorf("save photo: %w", err)
	}

	restaurant, oldURL, err := s.menegmentRepository.SetRestaurantPhoto(ctx, restID, url)
	if err != nil {
		core_files.RemovePhoto(&url)
		return domain.RestaurantResponse{}, fmt.Errorf("set restaurant photo in service: %w", err)
	}

	core_files.RemovePhoto(oldURL)

	return restaurant, nil
}
