package users_transport_http

import (
	"net/http"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_request "github.com/PassDelivery/pass-backend/internal/core/transport/http/request"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
)

func (s *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

	logger.Debug("invoke PatchUser handler")

	userID, err := core_jwt.GetIDFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "get id from context")
		return
	}

	var request UserPatchDTO

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "decode or validate error")
		return
	}

	user := userPatchFromRequest(request)

	userResponse, err := s.usersService.PatchUser(ctx, userID, user)
	if err != nil {
		responseHandler.ErrorResponse(err, "patch user")
		return
	}

	responseHandler.JSONResponse(userResponse, http.StatusOK)
}

func userPatchFromRequest(r UserPatchDTO) domain.UserPatch {
	return domain.UserPatch{
		Email: r.Email.ToDomain(),
		Name:  r.Name.ToDomain(),
		Phone: r.Phone.ToDomain(),
	}
}
