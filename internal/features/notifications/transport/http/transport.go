package notifications_transport_http

import (
	"context"
	"net/http"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	core_http_middleware "github.com/PassDelivery/pass-backend/internal/core/transport/http/middleware"
	core_http_server "github.com/PassDelivery/pass-backend/internal/core/transport/server"
	"github.com/google/uuid"
)

type NotificationsHTTPHandler struct {
	notificationsService NotificationsService
	token                *core_jwt.TokenService
}

type NotificationsService interface {
	GetNotifications(
		ctx context.Context,
	) ([]domain.NotificationResponse, error)
	ReadNotification(
		ctx context.Context,
		notificationID uuid.UUID,
	) (domain.NotificationResponse, error)
}

func NewNotificationsHTTPHandler(ns NotificationsService, token *core_jwt.TokenService) *NotificationsHTTPHandler {
	return &NotificationsHTTPHandler{
		notificationsService: ns,
		token:                token,
	}
}

func (h *NotificationsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/notifications",
			Handler: h.GetNotifications,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleCustomer,
					domain.RoleRestaurantOwner,
					domain.RoleCourier,
				),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/notifications/{id}/read",
			Handler: h.ReadNotification,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleCustomer,
					domain.RoleRestaurantOwner,
					domain.RoleCourier,
				),
			},
		},
	}
}
