package restaurants_transport_http

import (
	"context"
	"net/http"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_http_server "github.com/PassDelivery/pass-backend/internal/core/transport/server"
	"github.com/google/uuid"
)

type RestaurantsHTTPHandler struct {
	restaurantsService RestaurantsService
}

type RestaurantsService interface {
	GetRestaurants(
		ctx context.Context,
	) ([]domain.RestaurantCard, error)
	GetRestaurant(
		ctx context.Context,
		restID uuid.UUID,
	) (domain.RestaurantCard, error)
	GetMenu(
		ctx context.Context,
		restID uuid.UUID,
	) (domain.Menu, error)
}

func NewRestaurantsHTTPHandler(rs RestaurantsService) *RestaurantsHTTPHandler {
	return &RestaurantsHTTPHandler{
		restaurantsService: rs,
	}
}

func (h *RestaurantsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/restaurants",
			Handler: h.GetRestaurants,
		},
		{
			Method:  http.MethodGet,
			Path:    "/restaurants/{id}",
			Handler: h.GetRestaurant,
		},
		{
			Method:  http.MethodGet,
			Path:    "/restaurants/{id}/menu",
			Handler: h.GetMenu,
		},
	}
}
