package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_jwt "github.com/PassDelivery/pass-backend/internal/core/jwt"
	core_logger "github.com/PassDelivery/pass-backend/internal/core/logger"
	core_postgres_pool "github.com/PassDelivery/pass-backend/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/PassDelivery/pass-backend/internal/core/transport/http/middleware"
	core_http_server "github.com/PassDelivery/pass-backend/internal/core/transport/server"
	menegment_postgres_repository "github.com/PassDelivery/pass-backend/internal/features/menegment/repository/postgres"
	menegment_service "github.com/PassDelivery/pass-backend/internal/features/menegment/service"
	menegment_transport_http "github.com/PassDelivery/pass-backend/internal/features/menegment/transport/http"
	orders_postgres_repository "github.com/PassDelivery/pass-backend/internal/features/orders/repository/postgres"
	orders_service "github.com/PassDelivery/pass-backend/internal/features/orders/service"
	orders_transport_http "github.com/PassDelivery/pass-backend/internal/features/orders/transport/http"
	users_postgres_repository "github.com/PassDelivery/pass-backend/internal/features/users/repository/postgres"
	users_service "github.com/PassDelivery/pass-backend/internal/features/users/service"
	users_transport_http "github.com/PassDelivery/pass-backend/internal/features/users/transport/http"
	"go.uber.org/zap"

	restaurants_postgres_repository "github.com/PassDelivery/pass-backend/internal/features/restaurants/repository/postgres"
	restaurants_service "github.com/PassDelivery/pass-backend/internal/features/restaurants/service"
	restaurants_transport_http "github.com/PassDelivery/pass-backend/internal/features/restaurants/transport/http"

	cart_postgres_repository "github.com/PassDelivery/pass-backend/internal/features/cart/repository/postgres"
	cart_service "github.com/PassDelivery/pass-backend/internal/features/cart/service"
	cart_transport_http "github.com/PassDelivery/pass-backend/internal/features/cart/transport/http"
)

func main() {
	fmt.Println("hello world")

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	log, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("faild to init logger:", err)
		os.Exit(1)
	}
	defer log.Close()

	log.Debug("Initializing postgres connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		log.Fatal("Failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	token := core_jwt.NewTokenService(24 * time.Hour)

	log.Debug("Initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository, *token)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService, token)

	log.Debug("Initializing feature", zap.String("feature", "menagment-restaurant"))
	menegmentRepository := menegment_postgres_repository.NewMenegmentRepository(pool)
	menegmentService := menegment_service.NewMenegmentService(menegmentRepository)
	menegmentTransportHTTP := menegment_transport_http.NewMenegmentHTTPHandler(menegmentService, token)

	log.Debug("Initializing feature", zap.String("feature", "orders"))
	ordersRepository := orders_postgres_repository.NewOrdersRepository(pool)
	ordersService := orders_service.NewOrdersService(ordersRepository)
	ordersTransportHTTP := orders_transport_http.NewOrdersHTTPHandler(ordersService, token)


	log.Debug("Initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewHTTPConfigMust(),
		log,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(log),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
		core_http_middleware.CORS("http://localhost:5173"),
	)

	log.Debug("Initializing feature", zap.String("feature", "restaurants"))
	restaurantsRepository := restaurants_postgres_repository.NewRestaurantsRepository(pool)
	restaurantsService := restaurants_service.NewRestaurantsService(restaurantsRepository)
	restaurantsTransportHTTP := restaurants_transport_http.NewRestaurantsHTTPHandler(restaurantsService)

	log.Debug("Initializing feature", zap.String("feature", "cart"))
	cartRepository := cart_postgres_repository.NewCartRepository(pool)
	cartService := cart_service.NewCartService(cartRepository)
	cartTransportHTTP := cart_transport_http.NewCartHTTPTransport(cartService, token)

	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(restaurantsTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(menegmentTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(ordersTransportHTTP.Routes()...)
	apiVersionRouter.RegisterRoutes(cartTransportHTTP.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		log.Error("HTTP server run error", zap.Error(err))
	}
}
