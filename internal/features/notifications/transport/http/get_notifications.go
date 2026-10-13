package notifications_transport_http

import (
	"net/http"

	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
)

func (h *NotificationsHTTPHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke GetNotifications handler")

	notifications, err := h.notificationsService.GetNotifications(ctx)
	if err != nil {
		responseWriter.ErrorResponse(err, "get notifications")
		return
	}

	responseWriter.JSONResponse(notificationsToDTO(notifications), http.StatusOK)
}
