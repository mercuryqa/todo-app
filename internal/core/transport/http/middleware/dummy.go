package core_http_middleware

import (
	"fmt"
	"net/http"

	core_logger "github.com/mercuryqa/todo-app/internal/core/logger"
)

/*
	Работает только на GET users V1
	Либо на всех ручках V2
*/

// Dummy - добавляет requestID (id запроса) в хедер X-Request-ID
func Dummy(s string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)

			log.Debug(fmt.Sprintf("TEST -> before %s", s))

			next.ServeHTTP(w, r)

			log.Debug(fmt.Sprintf("TEST <- after %s", s))
		})
	}
}
