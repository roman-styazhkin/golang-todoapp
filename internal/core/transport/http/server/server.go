package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/roman-styazhkin/golang-todoapp/docs"
	core_logger "github.com/roman-styazhkin/golang-todoapp/internal/core/logger"
	core_http_middleware "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
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

func (s *HttpServer) RegisterSwagger() {
	s.mux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
			httpSwagger.DefaultModelsExpandDepth(-1),
		),
	)

	s.mux.HandleFunc(
		"/swagger/doc.json",
		func(rw http.ResponseWriter, r *http.Request) {
			r.Header.Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			_, _ = rw.Write([]byte(docs.SwaggerInfo.ReadDoc()))
		},
	)
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

func (s *HttpServer) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		s.mux.Handle(pattern, route.Handler)
	}
}
