package orders_transport_http

import (
	"context"
	"net/http"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	core_http_middleware "github.com/PassDelivery/pass-backend/internal/core/transport/http/middleware"
	core_http_server "github.com/PassDelivery/pass-backend/internal/core/transport/server"
	"github.com/google/uuid"
)

type OrdersHTTPHandler struct {
	ordersService OrdersService
	token         *core_jwt.TokenService
}

type OrdersService interface {
	ChangeOrderStatus(
		ctx context.Context,
		orderID uuid.UUID,
		status domain.OrderStatus,
	) (domain.OrderResponse, error)
}

func NewOrdersHTTPHandler(os OrdersService, token *core_jwt.TokenService) *OrdersHTTPHandler {
	return &OrdersHTTPHandler{
		ordersService: os,
		token:         token,
	}
}

func (h *OrdersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPatch,
			Path:    "/internal/orders/{id}/status",
			Handler: h.ChangeOrderStatus,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleRestaurantOwner,
				),
			},
		},
	}
}
