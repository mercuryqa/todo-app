// Package main содержит точку входа в приложение.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	coreLogger "github.com/mercuryqa/todo-app/internal/core/logger"
	core_pgx_pool "github.com/mercuryqa/todo-app/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/mercuryqa/todo-app/internal/core/transport/http/middleware"
	"github.com/mercuryqa/todo-app/internal/core/transport/http/server"
	users_postgres_repository "github.com/mercuryqa/todo-app/internal/features/users/repository/postgres"
	users_service "github.com/mercuryqa/todo-app/internal/features/users/service"
	users_transport_http "github.com/mercuryqa/todo-app/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	fmt.Println("Hello Todo app")

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	logger, err := coreLogger.NewLogger(coreLogger.NewCLoggerConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger", err)
		os.Exit(1)
	}
	defer logger.Close()

	// Создаём пулл соединений с PostgreSQL через библиотеку pgx.
	// Пул переиспользует соединения, что гораздо эффективнее,
	// чем открывать новое соединение на каждый SQL запрос.
	logger.Debug("initializing postgres connection pool")

	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("Initializing HTTP server")
	httpServer := server.NewHTTPServer(
		server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Recovery(),
	)

	apiVersionRouterV1 := server.NewApiVersionRouter(server.ApiVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)

	apiVersionRouterV2 := server.NewApiVersionRouter(
		server.ApiVersion2,
		core_http_middleware.Dummy("api v2 middleware"),
	)
	apiVersionRouterV2.RegisterRoutes(usersTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(
		apiVersionRouterV1,
		apiVersionRouterV2,
	)

	if err = httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
