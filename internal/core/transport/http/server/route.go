package server

import (
	"net/http"

	core_http_middleware "github.com/mercuryqa/todo-app/internal/core/transport/http/middleware"
)

// Route представляет HTTP-маршрут с методом, путём и обработчиком запроса.
type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []core_http_middleware.Middleware
}

// WithMiddleware применяет middleware маршрута к обработчику и возвращает готовый http.Handler.
func (r *Route) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r.Handler,
		r.Middleware...,
	)
}
