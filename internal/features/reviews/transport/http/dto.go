package reviews_transport_http

import (
	"fmt"
	"strings"
	"time"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	"github.com/google/uuid"
)

type ReviewRequestDTO struct {
	Rating        int     `json:"rating"`
	CourierRating *int    `json:"courier_rating"`
	Body          *string `json:"body"`
}

func (d *ReviewRequestDTO) Validate() error {
	if d.Rating < 1 || d.Rating > 5 {
		return fmt.Errorf("rating must be 1..5: %w", core_errors.ErrInvalidArgument)
	}
	if d.CourierRating != nil && (*d.CourierRating < 1 || *d.CourierRating > 5) {
		return fmt.Errorf("courier_rating must be 1..5: %w", core_errors.ErrInvalidArgument)
	}
	if d.Body != nil && len([]rune(*d.Body)) > 1000 {
		return fmt.Errorf("body len must be <= 1000: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

type ReviewResponseDTO struct {
	ID            uuid.UUID `json:"id"`
	OrderID       uuid.UUID `json:"order_id"`
	RestID        uuid.UUID `json:"restaurant_id"`
	AuthorName    string    `json:"author_name"`
	Rating        int       `json:"rating"`
	CourierRating *int      `json:"courier_rating,omitempty"`
	Body          *string   `json:"body"`
	CreatedAt     time.Time `json:"created_at"`
}

func reviewToDomain(review ReviewRequestDTO) domain.ReviewRequest {
	var body *string
	if review.Body != nil && strings.TrimSpace(*review.Body) != "" {
		trimmed := strings.TrimSpace(*review.Body)
		body = &trimmed
	}

	return domain.ReviewRequest{
		Rating:        review.Rating,
		CourierRating: review.CourierRating,
		Body:          body,
	}
}

func reviewToDTO(review domain.ReviewResponse) ReviewResponseDTO {
	return ReviewResponseDTO{
		ID:            review.ID,
		OrderID:       review.OrderID,
		RestID:        review.RestID,
		AuthorName:    review.AuthorName,
		Rating:        review.Rating,
		CourierRating: review.CourierRating,
		Body:          review.Body,
		CreatedAt:     review.CreatedAt,
	}
}

func reviewsToDTO(reviews []domain.ReviewResponse) []ReviewResponseDTO {
	res := make([]ReviewResponseDTO, len(reviews))
	for k, v := range reviews {
		res[k] = reviewToDTO(v)
	}

	return res
}
