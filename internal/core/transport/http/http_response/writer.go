package resonse

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// statusCodeUnitialized - устанавливаем значение -1 для статус кода, чтобы понимать, что он еще не был записан
var (
	statusCodeUnitialized = -1
)

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     statusCodeUnitialized,
	}
}

func (rw *ResponseWriter) WriteHeader(stausCode int) {
	rw.ResponseWriter.WriteHeader(stausCode)
	rw.statusCode = stausCode
}

// GetStatusCodeOrPanic - получаю значение приватного поля statusCode
func (rw *ResponseWriter) GetStatusCodeOrPanic() int {
	if rw.statusCode == statusCodeUnitialized {
		panic("no status code set")
	}
	return rw.statusCode
}
