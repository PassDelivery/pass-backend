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

func (h *MenegmentHTTPHandler) UploadRestaurantPhoto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke UploadRestaurantPhoto handler")

	id := r.PathValue("id")
	restID, err := uuid.Parse(id)
	if err != nil {
		responseWriter.ErrorResponse(fmt.Errorf("parse restaurant id: %w", core_errors.ErrInvalidArgument), "parse to uuid")
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

	restaurant, err := h.menegmentService.UploadRestaurantPhoto(ctx, restID, data)
	if err != nil {
		responseWriter.ErrorResponse(err, "upload restaurant photo")
		return
	}

	responseWriter.JSONResponse(DTOFromDomain(restaurant), http.StatusOK)
}
