package reviews_transport_http

import (
	"fmt"
	"net/http"

	core_errors "github.com/PassDelivery/pass-backend/internal/core/errors"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_request "github.com/PassDelivery/pass-backend/internal/core/transport/http/request"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
	"github.com/google/uuid"
)

func (h *ReviewsHTTPHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke CreateReview handler")

	id := r.PathValue("id")
	orderID, err := uuid.Parse(id)
	if err != nil {
		responseWriter.ErrorResponse(fmt.Errorf("parse order id: %w", core_errors.ErrInvalidArgument), "parse to uuid")
		return
	}

	var reviewReq ReviewRequestDTO

	if err := core_http_request.DecodeAndValidateRequest(r, &reviewReq); err != nil {
		responseWriter.ErrorResponse(err, "decode or validate")
		return
	}

	review, err := h.reviewsService.CreateReview(ctx, reviewToDomain(reviewReq), orderID)
	if err != nil {
		responseWriter.ErrorResponse(err, "create review")
		return
	}

	responseWriter.JSONResponse(reviewToDTO(review), http.StatusCreated)
}
