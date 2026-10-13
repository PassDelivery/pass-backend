package cart_transport_http

import (
	"context"
	"net/http"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	core_http_middleware "github.com/PassDelivery/pass-backend/internal/core/transport/http/middleware"
	core_http_server "github.com/PassDelivery/pass-backend/internal/core/transport/server"
	"github.com/google/uuid"
)

type CartHTTPTransport struct {
	cartService CartService
	token       *core_jwt.TokenService
}

type CartService interface {
	GetCart(
		ctx context.Context,
	) (domain.Cart, error)
	AddItem(
		ctx context.Context,
		item domain.CartItemRequest,
	) (domain.Cart, error)
	PatchItem(
		ctx context.Context,
		itemID uuid.UUID,
		quantity int,
	) (domain.Cart, error)
	DeleteItem(
		ctx context.Context,
		itemID uuid.UUID,
	) (domain.Cart, error)
	ClearCart(
		ctx context.Context,
	) error
}

func NewCartHTTPTransport(cs CartService, token *core_jwt.TokenService) *CartHTTPTransport {
	return &CartHTTPTransport{
		cartService: cs,
		token:       token,
	}
}

func (h *CartHTTPTransport) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/cart",
			Handler: h.GetCart,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleCustomer,
				),
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/cart",
			Handler: h.ClearCart,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleCustomer,
				),
			},
		},
		{
			Method:  http.MethodPost,
			Path:    "/cart/items",
			Handler: h.AddItem,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleCustomer,
				),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/cart/items/{id}",
			Handler: h.PatchItem,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleCustomer,
				),
			},
		},
		{
			Method:  http.MethodDelete,
			Path:    "/cart/items/{id}",
			Handler: h.DeleteItem,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleCustomer,
				),
			},
		},
	}
}
