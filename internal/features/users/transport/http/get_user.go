package users_transport_http

import (
	"net/http"

	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
)

func (h *UsersHTTPHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

	

	logger.Debug("invoke GetUser handler")

	user, err := h.usersService.GetUser(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "get user")
		return
	}

	userResp := GetUserResponseDTO{
		ID: user.ID,
		Email: user.Email,
		Name: user.Name,
		Phone: user.Phone,
		Role: string(user.Role),
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
	}

	responseHandler.JSONResponse(userResp, http.StatusOK)
}
