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

func (s *MenegmentService) UploadDishPhoto(
	ctx context.Context,
	dishID uuid.UUID,
	data []byte,
) (domain.DishResponse, error) {
	ownerID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.DishResponse{}, fmt.Errorf("get id from context: %w", err)
	}

	dish, err := s.menegmentRepository.GetDish(ctx, dishID)
	if err != nil {
		return domain.DishResponse{}, fmt.Errorf("get dish: %w", err)
	}

	isOwner, err := s.menegmentRepository.IsRestaurantOwner(ctx, dish.RestID, ownerID)
	if err != nil {
		return domain.DishResponse{}, fmt.Errorf("check restaurant owner: %w", err)
	}
	if !isOwner {
		return domain.DishResponse{}, fmt.Errorf("dish not found: %w", core_errors.ErrNotFound)
	}

	url, err := core_files.SavePhoto(core_files.PhotoKindDishes, data)
	if err != nil {
		return domain.DishResponse{}, fmt.Errorf("save photo: %w", err)
	}

	updated, err := s.menegmentRepository.SetDishPhoto(ctx, dishID, url)
	if err != nil {
		core_files.RemovePhoto(&url)
		return domain.DishResponse{}, fmt.Errorf("set dish photo in service: %w", err)
	}

	core_files.RemovePhoto(dish.PhotoURL)

	return updated, nil
}
