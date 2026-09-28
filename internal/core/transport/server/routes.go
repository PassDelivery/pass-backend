package core_http_server

import (
	"net/http"

	core_http_middleware "github.com/PassDelivery/pass-backend/internal/core/transport/http/middleware"
)

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
	Middleware []core_http_middleware.Middleware
}

// func NewRoute(
// 	method string,
// 	path string,
// 	handler http.HandlerFunc,
// 	middleware []core_http_middleware.Middleware,
// ) Route {
// 	return Route{
// 		Method:  method,
// 		Path:    path,
// 		Handler: handler,
// 	}
// }
