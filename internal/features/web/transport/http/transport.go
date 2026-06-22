package web_transport_http

import (
	core_http_server "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/server"
)

type WebHttpHandler struct {
	webService WebService
}

type WebService interface {
	GetMainPage() ([]byte, error)
}

func NewWebHttpHandler(webHttpService WebService) *WebHttpHandler {
	return &WebHttpHandler{
		webService: webHttpService,
	}
}

func (h *WebHttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Path:    "/",
			Handler: h.GetMainPage,
		},
	}
}
