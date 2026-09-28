package users_transport_http

import (
	"context"
	"net/http"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	core_http_middleware "github.com/PassDelivery/pass-backend/internal/core/transport/http/middleware"
	core_http_server "github.com/PassDelivery/pass-backend/internal/core/transport/server"
	"github.com/google/uuid"
)

type UsersHTTPHandler struct {
	usersService UsersService
	token        *core_jwt.TokenService
}

type UsersService interface {
	RegisterUser(
		ctx context.Context,
		user domain.UserRequest,
	) (domain.UserResponse, error)
	LoginUser(
		ctx context.Context,
		email string,
		password string,
	) (domain.UserResponse, error)
	GetUser(
		ctx context.Context,
	) (domain.GetUserResponse, error)
	PatchUser(
		ctx context.Context,
		id uuid.UUID,
		user domain.UserPatch,
	) (domain.GetUserResponse, error)
}

func NewUsersHTTPHandler(usersService UsersService, token *core_jwt.TokenService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
		token:        token,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/auth/register",
			Handler: h.RegisterUser,
		},
		{
			Method: http.MethodPost,
			Path: "/auth/login",
			Handler: h.LoginUser,
		},
		{
			Method:  http.MethodPost,
			Path:    "/internal/restaurants",
			Handler: func(w http.ResponseWriter, r *http.Request) {},
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleRestaurantOwner,
				),
			},
		},
		{
			Method: http.MethodGet,
			Path: "/users/me",
			Handler: h.GetUser,
			Middleware: []core_http_middleware.Middleware{
				core_http_middleware.Auth(h.token),
				core_http_middleware.RequireRole(
					domain.RoleCustomer,
				),
			},
		},
	}
}
