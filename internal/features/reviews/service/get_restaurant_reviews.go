package reviews_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
)

func (s *ReviewsService) GetRestaurantReviews(
	ctx context.Context,
	restID uuid.UUID,
) ([]domain.ReviewResponse, error) {
	exists, err := s.reviewsRepository.IsRestaurantExists(ctx, restID)
	if err != nil {
		return nil, fmt.Errorf("check restaurant: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("restaurant not found: %w", core_errors.ErrNotFound)
	}

	reviews, err := s.reviewsRepository.GetRestaurantReviews(ctx, restID)
	if err != nil {
		return nil, fmt.Errorf("get restaurant reviews in service: %w", err)
	}

	return reviews, nil
}
