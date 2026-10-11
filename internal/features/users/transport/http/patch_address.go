package users_transport_http

import (
	"net/http"

	"github.com/PassDelivery/pass-backend/internal/core/domain"
	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_request "github.com/PassDelivery/pass-backend/internal/core/transport/http/request"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *UsersHTTPHandler) PatchAddress(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke PatchAddress handler")

	id := r.PathValue("id")

	addrID, err := uuid.Parse(id)
	if err != nil {
		responseWriter.ErrorResponse(err, core_errors.ErrInvalidArgument.Error())
		return
	}

	var request AddressPatchDTO
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseWriter.ErrorResponse(err, "decode or validate request")
		return
	}

	addressDomain := addressPatchFromDTO(request)

	address, err := h.usersService.PatchAddress(ctx, addressDomain, addrID)
	if err != nil {
		responseWriter.ErrorResponse(err, "failed to patch address")
		return
	}

	responseWriter.JSONResponse(addressResponseDTOFromDomain(address), http.StatusOK)
}

func addressPatchFromDTO(p AddressPatchDTO) domain.AddressPatch {
	return domain.AddressPatch{
		City:      p.City.ToDomain(),
		Street:    p.Street.ToDomain(),
		House:     p.House.ToDomain(),
		Apartment: p.Apartment.ToDomain(),
		Comment:   p.Comment.ToDomain(),
		IsDefault: p.IsDefault.ToDomain(),
	}
}
