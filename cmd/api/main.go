package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_logger "github.com/roman-styazhkin/golang-todoapp/internal/core/logger"
	core_pgx_pool "github.com/roman-styazhkin/golang-todoapp/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/server"
	tasks_repository "github.com/roman-styazhkin/golang-todoapp/internal/features/tasks/repository/postgres"
	tasks_service "github.com/roman-styazhkin/golang-todoapp/internal/features/tasks/service"
	tasks_transport_http "github.com/roman-styazhkin/golang-todoapp/internal/features/tasks/transport/http"
	users_repository "github.com/roman-styazhkin/golang-todoapp/internal/features/users/repository/postgres"
	users_service "github.com/roman-styazhkin/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/roman-styazhkin/golang-todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

var (
	timeZone = time.UTC
)

func main() {
	time.Local = timeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init logger")
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application timeZone:", zap.Any("timeZone: ", timeZone))

	logger.Debug("init pool...")
	pool, err := core_pgx_pool.NewPool(ctx, core_pgx_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to create pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("users feature...")
	usersRepository := users_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersHandler := users_transport_http.NewUsersHttpHandler(usersService)

	logger.Debug("tasks feature...")
	tasksRepository := tasks_repository.NewTaskRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksHandler := tasks_transport_http.NewTasksHttpHandler(tasksService)

	logger.Debug("init api version router...")
	apiVersionRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersionRouter1)
	apiVersionRouter.RegisterRoutes(usersHandler.Routes()...)
	apiVersionRouter.RegisterRoutes(tasksHandler.Routes()...)

	logger.Debug("init server...")
	httpServer := core_http_server.NewHttpServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	httpServer.RegisterRouters(apiVersionRouter)

	logger.Debug("start router...")
	if err = httpServer.Run(ctx); err != nil {
		logger.Error("failed to start http server", zap.Error(err))
	}
}
