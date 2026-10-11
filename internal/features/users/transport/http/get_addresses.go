package users_transport_http

import (
	"net/http"

	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
)

func (h *UsersHTTPHandler) GetAddresses(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke GetAddresses handler")

	addresses, err := h.usersService.GetAddresses(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "get addresses failed")
		return
	}
	address := NewAddressesResponseDTO(addresses)

	responseHandler.JSONResponse(address, http.StatusOK)
}
