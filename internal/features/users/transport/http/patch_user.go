package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_request "github.com/PassDelivery/pass-backend/internal/core/transport/http/request"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
)

func (d *UserPatchDTO) Validate() error {
	if d.Email.Set {
		if d.Email.Value == nil {
			return fmt.Errorf("Email cant be null")
		}
		if ok := strings.Contains(*d.Email.Value, "@"); !ok {
			return fmt.Errorf("invalid symbol in Email")
		}
		if len(*d.Email.Value) < 5 {
			return fmt.Errorf("Email symbols count cant be < 5")
		}
	}

	if d.Name.Set {
		if d.Name.Value == nil {
			return fmt.Errorf("Name cant be null")
		}
		if len(*d.Name.Value) < 2 {
			return fmt.Errorf("Name symbols count cant be < 2")
		}
	}

	if d.Phone.Set {
		if d.Phone.Value == nil {
			return fmt.Errorf("Phone cant be null")
		}
		if ok := strings.Contains(*d.Phone.Value, "+"); !ok {
			return fmt.Errorf("invalid symbol in Email")
		}
		if len(*d.Email.Value) < 10 {
			return fmt.Errorf("Email symbols count cant be < 10")
		}
		if len(*d.Email.Value) > 15 {
			return fmt.Errorf("Email symbols count cant be > 15")
		}
	}

	return nil
}

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
	}

	responseHandler.JSONResponse(userResponse, http.StatusOK)
}

func userPatchFromRequest(r UserPatchDTO) domain.UserPatch {
	return domain.UserPatch{
		Email: r.Email.ToDomain(),
		Name: r.Name.ToDomain(),
		Phone: r.Phone.ToDomain(),
	}
}
