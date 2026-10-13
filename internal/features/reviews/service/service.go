package reviews_service

import (
	"context"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	"github.com/google/uuid"
)

type ReviewsService struct {
	reviewsRepository ReviewsRepository
}

type ReviewsRepository interface {
	GetReviewOrder(
		ctx context.Context,
		orderID uuid.UUID,
	) (domain.ReviewOrder, error)
	CreateReview(
		ctx context.Context,
		review domain.ReviewRequest,
		orderID uuid.UUID,
		userID uuid.UUID,
		restID uuid.UUID,
	) (domain.ReviewResponse, error)
	IsRestaurantExists(
		ctx context.Context,
		restID uuid.UUID,
	) (bool, error)
	GetRestaurantReviews(
		ctx context.Context,
		restID uuid.UUID,
	) ([]domain.ReviewResponse, error)
}

func NewReviewsService(rr ReviewsRepository) *ReviewsService {
	return &ReviewsService{
		reviewsRepository: rr,
	}
}
