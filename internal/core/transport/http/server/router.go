// Package server содержит компоненты для запуска и настройки HTTP-сервера.
package server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/mercuryqa/todo-app/internal/core/transport/http/middleware"
)

type ApiVersion string

// ApiVersion содержит доступные версии API.
var (
	ApiVersion1 = ApiVersion("v1")
	ApiVersion2 = ApiVersion("v2")
	ApiVersion3 = ApiVersion("v3")
)

// APIVersionRouter представляет HTTP-маршрутизатор с поддержкой версионирования API.
type APIVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
	middleware []core_http_middleware.Middleware
}

// NewApiVersionRouter создаёт HTTP-маршрутизатор с указанной версией API.
func NewApiVersionRouter(
	apiVersion ApiVersion,
	middleware ...core_http_middleware.Middleware,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
		middleware: middleware,
	}
}

// RegisterRoutes регистрирует переданные HTTP-маршруты в маршрутизаторе.
func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.WithMiddleware())
	}
}

func (r *APIVersionRouter) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r,
		r.middleware...,
	)
}
