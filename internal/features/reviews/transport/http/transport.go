package reviews_transport_http

import (
	"context"
	"net/http"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	core_http_middleware "github.com/PassDelivery/pass-backend/internal/core/transport/http/middleware"
	core_http_server "github.com/PassDelivery/pass-backend/internal/core/transport/server"
	"github.com/google/uuid"
)

type ReviewsHTTPHandler struct {
	reviewsService ReviewsService
	token          *core_jwt.TokenService
}

type ReviewsService interface {
	CreateReview(
		ctx context.Context,
		review domain.ReviewRequest,
		orderID uuid.UUID,
	) (domain.ReviewResponse, error)
	GetRestaurantReviews(
		ctx context.Context,
		restID uuid.UUID,
	) ([]domain.ReviewResponse, error)
}

func NewReviewsHTTPHandler(rs ReviewsService, token *core_jwt.TokenService) *ReviewsHTTPHandler {
	return &ReviewsHTTPHandler{
		reviewsService: rs,
		token:          token,
	}
}

func (h *ReviewsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/orders/{id}/review",
			Handler: h.CreateReview,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleCustomer,
				),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/restaurants/{id}/reviews",
			Handler: h.GetRestaurantReviews,
		},
	}
}
