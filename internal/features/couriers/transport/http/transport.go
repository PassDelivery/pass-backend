package couriers_transport_http

import (
	"context"
	"net/http"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	core_http_middleware "github.com/PassDelivery/pass-backend/internal/core/transport/http/middleware"
	core_http_server "github.com/PassDelivery/pass-backend/internal/core/transport/server"
	"github.com/google/uuid"
)

type CouriersHTTPHandler struct {
	couriersService CouriersService
	token           *core_jwt.TokenService
}

type CouriersService interface {
	GetStatus(
		ctx context.Context,
	) (domain.Courier, error)
	SetStatus(
		ctx context.Context,
		status domain.CourierStatus,
	) (domain.Courier, error)
	GetDeliveries(
		ctx context.Context,
	) ([]domain.Delivery, error)
	AcceptDelivery(
		ctx context.Context,
		deliveryID uuid.UUID,
	) (domain.Delivery, error)
	RefuseDelivery(
		ctx context.Context,
		deliveryID uuid.UUID,
	) error
	ChangeDeliveryStatus(
		ctx context.Context,
		deliveryID uuid.UUID,
		status domain.DeliveryStatus,
	) (domain.Delivery, error)
}

func NewCouriersHTTPHandler(cs CouriersService, token *core_jwt.TokenService) *CouriersHTTPHandler {
	return &CouriersHTTPHandler{
		couriersService: cs,
		token:           token,
	}
}

func (h *CouriersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/courier/status",
			Handler: h.GetStatus,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleCourier,
				),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/courier/status",
			Handler: h.SetStatus,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleCourier,
				),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/courier/deliveries",
			Handler: h.GetDeliveries,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleCourier,
				),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/courier/deliveries/{id}/accept",
			Handler: h.AcceptDelivery,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleCourier,
				),
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/courier/deliveries/{id}/refuse",
			Handler: h.RefuseDelivery,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleCourier,
				),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/courier/deliveries/{id}/status",
			Handler: h.ChangeDeliveryStatus,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleCourier,
				),
			},
		},
	}
}
