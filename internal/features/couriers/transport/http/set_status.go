package couriers_transport_http

import (
	"net/http"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_request "github.com/PassDelivery/pass-backend/internal/core/transport/http/request"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
)

func (h *CouriersHTTPHandler) SetStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke SetStatus handler")

	var statusReq CourierStatusRequestDTO

	if err := core_http_request.DecodeAndValidateRequest(r, &statusReq); err != nil {
		responseWriter.ErrorResponse(err, "decode or validate")
		return
	}

	courier, err := h.couriersService.SetStatus(ctx, domain.CourierStatus(statusReq.Status))
	if err != nil {
		responseWriter.ErrorResponse(err, "set courier status")
		return
	}

	responseWriter.JSONResponse(courierToDTO(courier), http.StatusOK)
}
