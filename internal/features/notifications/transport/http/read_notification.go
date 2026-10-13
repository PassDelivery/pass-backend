package notifications_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *NotificationsHTTPHandler) ReadNotification(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke ReadNotification handler")

	id := r.PathValue("id")
	notificationID, err := uuid.Parse(id)
	if err != nil {
		responseWriter.ErrorResponse(fmt.Errorf("parse notification id: %w", core_errors.ErrInvalidArgument), "parse to uuid")
		return
	}

	notification, err := h.notificationsService.ReadNotification(ctx, notificationID)
	if err != nil {
		responseWriter.ErrorResponse(err, "read notification")
		return
	}

	responseWriter.JSONResponse(notificationToDTO(notification), http.StatusOK)
}
