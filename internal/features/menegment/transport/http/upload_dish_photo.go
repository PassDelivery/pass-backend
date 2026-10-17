package menegment_transport_http

import (
	"fmt"
	"io"
	"net/http"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_files "github.com/PassDelivery/pass-backend/internal/core/files"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *MenegmentHTTPHandler) UploadDishPhoto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke UploadDishPhoto handler")

	id := r.PathValue("id")
	dishID, err := uuid.Parse(id)
	if err != nil {
		responseWriter.ErrorResponse(fmt.Errorf("parse dish id: %w", core_errors.ErrInvalidArgument), "parse to uuid")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, core_files.MaxPhotoSize+1<<20)
	file, _, err := r.FormFile("photo")
	if err != nil {
		responseWriter.ErrorResponse(fmt.Errorf("read photo from form: %v: %w", err, core_errors.ErrInvalidArgument), "read photo")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, core_files.MaxPhotoSize+1))
	if err != nil {
		responseWriter.ErrorResponse(fmt.Errorf("read photo: %v: %w", err, core_errors.ErrInvalidArgument), "read photo")
		return
	}

	dish, err := h.menegmentService.UploadDishPhoto(ctx, dishID, data)
	if err != nil {
		responseWriter.ErrorResponse(err, "upload dish photo")
		return
	}

	responseWriter.JSONResponse(dishToDTO(dish), http.StatusOK)
}
