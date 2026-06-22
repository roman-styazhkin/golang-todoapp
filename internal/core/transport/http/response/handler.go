package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/roman-styazhkin/golang-todoapp/internal/core/errors"
	core_logger "github.com/roman-styazhkin/golang-todoapp/internal/core/logger"
	"go.uber.org/zap"
)

type HttpResponseHandler struct {
	log *core_logger.Logger
	rw  http.ResponseWriter
}

func NewHttpResponseHandler(
	rw http.ResponseWriter,
	log *core_logger.Logger,
) *HttpResponseHandler {
	return &HttpResponseHandler{
		rw:  rw,
		log: log,
	}
}

func (h *HttpResponseHandler) PanicResponse(p any, message string) {
	errMessage := fmt.Errorf("unexpected panic, %v", p)

	h.log.Error(message, zap.Error(errMessage))
	h.errResponse(errMessage, message, http.StatusInternalServerError)
}

func (h *HttpResponseHandler) ErrorResponse(
	message string,
	err error,
) {
	var (
		statusCode int
		logFunc    func(message string, fields ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn
	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug
	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(message, zap.Error(err))
	h.errResponse(err, message, statusCode)
}

func (h *HttpResponseHandler) errResponse(
	err error,
	message string,
	statusCode int,
) {
	response := ErrResponse{
		Err:     err.Error(),
		Message: message,
	}

	h.JSONResponse(response, statusCode)
}

func (h *HttpResponseHandler) JSONResponse(
	response any,
	statusCode int,
) {
	h.rw.Header().Set("Content-Type", "application/json")
	h.rw.WriteHeader(statusCode)
	if err := json.NewEncoder(h.rw).Encode(response); err != nil {
		h.log.Error("failed to encode response", zap.Error(err))
	}
}

func (h *HttpResponseHandler) NoContentResponse() {
	h.rw.WriteHeader(http.StatusNoContent)
}

func (h *HttpResponseHandler) HtmlResponse(html []byte) {
	h.rw.WriteHeader(http.StatusOK)
	h.rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := h.rw.Write(html); err != nil {
		h.log.Error("write html http response", zap.Error(err))
	}
}
