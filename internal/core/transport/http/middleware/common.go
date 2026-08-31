package core_http_middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	coreLogger "github.com/mercuryqa/todo-app/internal/core/logger"
	"github.com/mercuryqa/todo-app/internal/core/transport/http/http_response"
	"go.uber.org/zap"
)

// type ctxKey string // 1/4

const (
	requestIDHeader = "X-Request-ID"
	// requestIDKey ctxKey = "request_id" // 2/4
)

// RequestID - добавляет requestID (id запроса) в хедер X-Request-ID
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			//ctx := context.WithValue(r.Context(), requestIDKey, requestID) // 3/4

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
			//next.ServeHTTP(w, r.WithContext(ctx)) // 4/4
		})
	}
}

// Logger - добавляет request_id и url запроса в контекст
func Logger(log *coreLogger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := coreLogger.ToContext(r.Context(), l)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Trace - логирует входящие запросы и исходящие ответы
func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Логируем REQUEST
			ctx := r.Context()
			log := coreLogger.FromContext(ctx)
			// Оборачиваю responseWriter, чтобы сохранить статус код
			rw := core_http_response.NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				">>> incoming HTTP request",
				zap.String("http_method", r.Method),
				zap.Time("time", before.UTC()),
			)

			// передаю кастомный responseWriter
			next.ServeHTTP(rw, r)

			// Логируем RESPONSE

			// Чтобы получить статус код из запроса, нужна обертка которая будет запоминмть статус код

			log.Debug(
				">>> done HTTP request",
				// получаю записанный статус код
				zap.Int("status_code", rw.GetStatusCode()),
				zap.Duration("latency", time.Now().Sub(before)),
			)

		})
	}
}

// В первых двух middleware паника не отлавливается и приложение упадет
// Дальше везде будет отлавливаться паника

// Recovery - Отлавливает панику в http запросах (после выполнения запроса в defer)
func Recovery() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := coreLogger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)
			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(p, "during handle HTTP request got unexpected panic")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
