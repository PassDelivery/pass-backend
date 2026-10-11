package menegment_transport_http

import (
	"context"
	"net/http"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	core_http_middleware "github.com/PassDelivery/pass-backend/internal/core/transport/http/middleware"
	core_http_server "github.com/PassDelivery/pass-backend/internal/core/transport/server"
	"github.com/google/uuid"
)

type MenegmentHTTPHandler struct {
	menegmentService MenegmentService
	token            *core_jwt.TokenService
}

type MenegmentService interface {
	RegisterRestaurant(
		ctx context.Context,
		restaurant domain.RestaurantRequest,
	) (domain.RestaurantResponse, error)
	PatchRestaurant(
		ctx context.Context,
		rest domain.RestaurantPatch,
		restID uuid.UUID,
	) (domain.RestaurantResponse, error)
	AddCategory(
		ctx context.Context,
		category domain.CategoryRequest,
		restID uuid.UUID,
	) (domain.CategoryResponse, error)
	AddDish(
		ctx context.Context,
		dish domain.DishRequest,
		restID uuid.UUID,
	) (domain.DishResponse, error)
	PatchDish(
		ctx context.Context,
		dish domain.DishPatch,
		dishID uuid.UUID,
	) (domain.DishResponse, error)
	HideDish(
		ctx context.Context,
		dishID uuid.UUID,
	) error
}

func NewMenegmentHTTPHandler(ms MenegmentService, token *core_jwt.TokenService) *MenegmentHTTPHandler {
	return &MenegmentHTTPHandler{
		menegmentService: ms,
		token:            token,
	}
}

func (h *MenegmentHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/internal/restaurants",
			Handler: h.RegisterRestaurant,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleRestaurantOwner,
				),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/internal/restaurants/{id}",
			Handler: h.PatchRestaurant,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleRestaurantOwner,
				),
			},
		},
				{
			Method:  http.MethodPost,
			Path:    "/internal/restaurants/{id}/categories",
			Handler: h.AddCategory,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleRestaurantOwner,
				),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/internal/restaurants/{id}/dishes",
			Handler: h.AddDish,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleRestaurantOwner,
				),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/internal/dishes/{id}",
			Handler: h.PatchDish,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleRestaurantOwner,
				),
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/internal/dishes/{id}",
			Handler: h.HideDish,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleRestaurantOwner,
				),
			},
		},
	}
}
