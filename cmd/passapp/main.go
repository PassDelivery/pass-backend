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
	users_postgres_repository "github.com/PassDelivery/pass-backend/internal/features/users/repository/postgres"
	users_service "github.com/PassDelivery/pass-backend/internal/features/users/service"
	users_transport_http "github.com/PassDelivery/pass-backend/internal/features/users/transport/http"
	"go.uber.org/zap"
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

	log.Debug("Initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	token := core_jwt.NewTokenService(24 * time.Hour)
	usersService := users_service.NewUsersService(usersRepository, *token)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService, token)

	log.Debug("Initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewHTTPConfigMust(),
		log,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(log),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(),
	)

	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		log.Error("HTTP server run error", zap.Error(err))
	}
}
