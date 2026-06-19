package core_http_response

import "net/http"

var (
	StatusCodeUninitialized = -1
)

type HttpResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewHttpResponseWriter(rw http.ResponseWriter) *HttpResponseWriter {
	return &HttpResponseWriter{
		ResponseWriter: rw,
	}
}

func (rw *HttpResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func (rw *HttpResponseWriter) GetStatusCode() int {
	if rw.statusCode == StatusCodeUninitialized {
		return http.StatusOK
	}

	return rw.statusCode
}
