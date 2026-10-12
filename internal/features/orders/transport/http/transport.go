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
	CreateOrder(
		ctx context.Context,
		order domain.CreateOrderRequest,
	) (domain.OrderResponse, error)
	GetOrders(
		ctx context.Context,
	) ([]domain.OrderResponse, error)
	GetOrder(
		ctx context.Context,
		orderID uuid.UUID,
	) (domain.OrderResponse, error)
	CancelOrder(
		ctx context.Context,
		orderID uuid.UUID,
	) (domain.OrderResponse, error)
	GetStatusHistory(
		ctx context.Context,
		orderID uuid.UUID,
	) ([]domain.OrderStatusChange, error)
	GetRestaurantOrders(
		ctx context.Context,
		restID uuid.UUID,
	) ([]domain.OrderResponse, error)
}

func NewOrdersHTTPHandler(os OrdersService, token *core_jwt.TokenService) *OrdersHTTPHandler {
	return &OrdersHTTPHandler{
		ordersService: os,
		token:         token,
	}
}

func (h *OrdersHTTPHandler) Routes() []core_http_server.Route {
	customer := []core_http_middleware.Middleware{
		core_http_middleware.Auth(h.token),
		core_http_middleware.RequireRole(
			domain.RoleCustomer,
		),
	}

	return []core_http_server.Route{
		{
			Method:     http.MethodGet,
			Path:       "/orders",
			Handler:    h.GetOrders,
			Middleware: customer,
		},
		{
			Method:     http.MethodPost,
			Path:       "/orders",
			Handler:    h.CreateOrder,
			Middleware: customer,
		},
		{
			Method:     http.MethodGet,
			Path:       "/orders/{id}",
			Handler:    h.GetOrder,
			Middleware: customer,
		},
		{
			Method:     http.MethodPost,
			Path:       "/orders/{id}/cancel",
			Handler:    h.CancelOrder,
			Middleware: customer,
		},
		{
			Method:     http.MethodGet,
			Path:       "/orders/{id}/status-history",
			Handler:    h.GetStatusHistory,
			Middleware: customer,
		},
		{
			Method:  http.MethodGet,
			Path:    "/internal/restaurants/{id}/orders",
			Handler: h.GetRestaurantOrders,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleRestaurantOwner,
				),
			},
		},
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
