package users_transport_http

import (
	"fmt"
	"io"
	"net/http"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_files "github.com/PassDelivery/pass-backend/internal/core/files"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
)

func (h *UsersHTTPHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

	logger.Debug("invoke UploadAvatar handler")

	r.Body = http.MaxBytesReader(w, r.Body, core_files.MaxPhotoSize+1<<20)
	file, _, err := r.FormFile("photo")
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("read photo from form: %v: %w", err, core_errors.ErrInvalidArgument), "read photo")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, core_files.MaxPhotoSize+1))
	if err != nil {
		responseHandler.ErrorResponse(fmt.Errorf("read photo: %v: %w", err, core_errors.ErrInvalidArgument), "read photo")
		return
	}

	user, err := h.usersService.UploadAvatar(ctx, data)
	if err != nil {
		responseHandler.ErrorResponse(err, "upload avatar")
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
