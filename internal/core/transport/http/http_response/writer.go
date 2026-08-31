// Package core_http_response предоставляет инструменты для работы с HTTP-ответами.
package core_http_response

import "net/http"

// ResponseWriter расширяет http.ResponseWriter и хранит установленный HTTP-код ответа.
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// statusCodeUnitialized - устанавливаем значение -1 для статус кода, чтобы понимать, что он еще не был записан
var (
	statusCodeUnitialized = -1
)

// NewResponseWriter создаёт ResponseWriter и устанавливает начальное состояние HTTP-кода ответа.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     statusCodeUnitialized,
	}
}

// WriteHeader устанавливает HTTP-код ответа и передаёт его исходному ResponseWriter.
func (rw *ResponseWriter) WriteHeader(stausCode int) {
	rw.ResponseWriter.WriteHeader(stausCode)
	rw.statusCode = stausCode
}

// GetStatusCode возвращает HTTP-код ответа или вызывает panic, если код ответа ещё не установлен.
func (rw *ResponseWriter) GetStatusCode() int {
	if rw.statusCode == statusCodeUnitialized {
		return http.StatusOK
	}
	return rw.statusCode
}
