package core_http_server

import (
	"fmt"
	"net/http"
)

type ApiVersion string

var (
	ApiVersionRouter1 = ApiVersion("v1")
	ApiVersionRouter2 = ApiVersion("v2")
	ApiVersionRouter3 = ApiVersion("v3")
)

type ApiVersionRouter struct {
	mux        *http.ServeMux
	apiVersion ApiVersion
}

func NewApiVersionRouter(apiVersion ApiVersion) *ApiVersionRouter {
	return &ApiVersionRouter{
		mux:        http.NewServeMux(),
		apiVersion: apiVersion,
	}
}

func (r *ApiVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.mux.Handle(pattern, route.Handler)
	}
}
