package core_http_middleware

import "net/http"

// Middleware представляет middleware, оборачивающий HTTP-обработчик.
type Middleware func(handler http.Handler) http.Handler

// ChainMiddleware последовательно применяет переданные middleware к HTTP-обработчику.
// Middleware применяются в порядке, в котором они переданы в функцию.
func ChainMiddleware(h http.Handler, m ...Middleware) http.Handler {
	if len(m) == 0 {
		return h
	}

	// от последней к первой (чтобы первая вызывалась первой иначе будет наоборот)
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}

	return h
}
