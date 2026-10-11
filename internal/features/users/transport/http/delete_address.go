package users_transport_http

import (
	"net/http"

	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *UsersHTTPHandler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke DeleteAddress handler")

	id := r.PathValue("id")
	addrID, err := uuid.Parse(id)
	if err != nil {
		responseWriter.ErrorResponse(err, "parse to uuid")
		return
	}

	if err := h.usersService.DeleteAddress(ctx, addrID); err != nil {
		responseWriter.ErrorResponse(err, "delete address")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
