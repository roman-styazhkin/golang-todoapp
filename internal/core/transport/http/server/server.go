package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/roman-styazhkin/golang-todoapp/internal/core/logger"
	core_http_middleware "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

type HttpServer struct {
	mux        *http.ServeMux
	config     Config
	log        *core_logger.Logger
	middleware []core_http_middleware.Middleware
}

func NewHttpServer(
	config Config,
	log *core_logger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HttpServer {
	return &HttpServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (s *HttpServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)
		s.log.Debug("listen http server, ", zap.String("port:", s.config.Addr))
		err := server.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("failed to listen http server, %w", err)
		}
	case <-ctx.Done():
		s.log.Debug("shutdown http server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			s.log.Warn("failed to shutdown http server", zap.Error(err))
			return fmt.Errorf("failed to shutdown http server, %w", err)
		}

		s.log.Debug("shutdown http server success!")
	}

	return nil
}

func (s *HttpServer) RegisterRouters(routers ...*ApiVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)
		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router.mux),
		)
	}
}
