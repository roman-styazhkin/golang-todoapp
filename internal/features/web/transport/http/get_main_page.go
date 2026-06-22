package web_transport_http

import (
	"net/http"

	core_logger "github.com/roman-styazhkin/golang-todoapp/internal/core/logger"
	core_http_response "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/response"
)

func (h *WebHttpHandler) GetMainPage(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(rw, log)

	htmlFile, err := h.webService.GetMainPage()
	if err != nil {
		responseHandler.ErrorResponse("failed to get htmlFile", err)
		return
	}

	responseHandler.HtmlResponse(htmlFile)
}
