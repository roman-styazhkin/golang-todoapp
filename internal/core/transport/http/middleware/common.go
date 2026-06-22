package core_http_middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	core_logger "github.com/roman-styazhkin/golang-todoapp/internal/core/logger"
	core_http_response "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/response"
	"go.uber.org/zap"
)

const requestID = "X-Request-ID"

func CORS(allowedOriginsList []string) Middleware {
	allowedOriginsMap := make(map[string]struct{})

	for _, origin := range allowedOriginsList {
		allowedOriginsMap[origin] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if _, ok := allowedOriginsMap[origin]; ok {
				rw.Header().Set("Access-Control-Allow-Origin", origin)
				rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				rw.Header().Set("Access-Control-Allow-Headers", "Content-type, Authorization")
				rw.Header().Set("Vary", "Origin")
			}

			if r.Method == http.MethodOptions {
				rw.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			requestHeader := r.Header.Get(requestID)
			if requestHeader == "" {
				requestHeader = uuid.NewString()
			}

			rw.Header().Set(requestID, requestHeader)
			r.Header.Set(requestID, requestHeader)
			next.ServeHTTP(rw, r)
		})
	}
}

func Logger(logger *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			requestHeader := r.Header.Get(requestID)

			log := logger.With(
				zap.String("request_id", requestHeader),
				zap.String("url", r.URL.String()),
			)

			ctx := core_logger.ToContext(r.Context(), log)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())
			httpResponseWriter := core_http_response.NewHttpResponseWriter(rw)

			start := time.Now()

			log.Debug(
				"<--- incoming http request",
				zap.Time("start:", start.UTC()),
			)

			next.ServeHTTP(httpResponseWriter, r)

			log.Debug(
				"<--- done http request",
				zap.Int("status_code", httpResponseWriter.GetStatusCode()),
				zap.String("request_method", r.Method),
				zap.Duration("latency:", time.Now().Sub(start)),
			)
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())
			responseHttpHandler := core_http_response.NewHttpResponseHandler(rw, log)

			defer func() {
				if p := recover(); p != nil {
					responseHttpHandler.PanicResponse(p, "incoming http request panic")
				}
			}()

			next.ServeHTTP(rw, r)
		})
	}
}
