package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReviewRequest struct {
	Rating int
	Body   *string
}

type ReviewResponse struct {
	ID         uuid.UUID
	OrderID    uuid.UUID
	RestID     uuid.UUID
	AuthorName string
	Rating     int
	Body       *string
	CreatedAt  time.Time
}

type ReviewOrder struct {
	UserID uuid.UUID
	RestID uuid.UUID
	Status OrderStatus
}
