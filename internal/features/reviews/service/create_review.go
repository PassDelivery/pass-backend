package reviews_service

import (
	"context"
	"fmt"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	"github.com/google/uuid"
)

func (s *ReviewsService) CreateReview(
	ctx context.Context,
	review domain.ReviewRequest,
	orderID uuid.UUID,
) (domain.ReviewResponse, error) {
	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		return domain.ReviewResponse{}, fmt.Errorf("get id from context: %w", err)
	}

	order, err := s.reviewsRepository.GetReviewOrder(ctx, orderID)
	if err != nil {
		return domain.ReviewResponse{}, fmt.Errorf("get order: %w", err)
	}

	if order.UserID != userID {
		return domain.ReviewResponse{}, fmt.Errorf("order not found: %w", core_errors.ErrNotFound)
	}

	if order.Status != domain.OrderStatusDelivered {
		return domain.ReviewResponse{}, fmt.Errorf("order is %s, review only after delivery: %w", order.Status, core_errors.ErrConflict)
	}

	reviewResp, err := s.reviewsRepository.CreateReview(ctx, review, orderID, userID, order.RestID)
	if err != nil {
		return domain.ReviewResponse{}, fmt.Errorf("create review in service: %w", err)
	}

	return reviewResp, nil
}
