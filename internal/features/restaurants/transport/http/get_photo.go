package restaurants_transport_http

import (
	"net/http"

	core_files "github.com/PassDelivery/pass-backend/internal/core/files"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_http_response "github.com/PassDelivery/pass-backend/internal/core/transport/http/response"
)

func (h *RestaurantsHTTPHandler) GetPhoto(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseWriter := core_http_response.NewHTTPResponseHandler(log, w)

	log.Debug("invoke GetPhoto handler")

	path, err := core_files.PhotoPath(r.PathValue("kind"), r.PathValue("name"))
	if err != nil {
		responseWriter.ErrorResponse(err, "get photo path")
		return
	}

	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	http.ServeFile(w, r, path)
}
