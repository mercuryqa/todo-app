package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/mercuryqa/todo-app/internal/core/errors"
	coreLogger "github.com/mercuryqa/todo-app/internal/core/logger"
	"go.uber.org/zap"
)

// HTTPResponseHandler обрабатывает HTTP-ответы и предоставляет доступ к логгеру.
type HTTPResponseHandler struct {
	log *coreLogger.Logger
	rw  http.ResponseWriter
}

// NewHTTPResponseHandler создаёт обработчик HTTP-ответов с переданным логгером и ResponseWriter.
func NewHTTPResponseHandler(
	log *coreLogger.Logger,
	rw http.ResponseWriter,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw:  rw,
	}
}

// JSONResponse сериализует responseBody в JSON и записывает в ответ с указанным статус-кодом.
// Content-Type автоматически определяется json.NewEncoder.
func (h *HTTPResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

// NoContentResponse отправляет HTTP 204 No Content — используется при успешном DELETE.
func (h *HTTPResponseHandler) NoContentResponse() {
	h.rw.WriteHeader(http.StatusNoContent)
}

// ErrorResponse транслирует core ошибку в HTTP-статус через errors.Is().
func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	// не совсем ошибка, но нештатная ситуация (почему делают такие запросы?) - warn
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	// not found - проблема на стороне клиента - debug
	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug

	// не совсем ошибка, но нештатная ситуация (почему делают такие запросы?) - warn
	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))
	h.errorResponse(
		statusCode,
		err,
		msg,
	)
}

// PanicResponse обрабатывает панику, возникшую при обработке HTTP-запроса,
// и возвращает клиенту ответ с кодом 500 и информацией об ошибке.
func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))
	h.errorResponse(
		statusCode,
		err,
		msg,
	)

}

// errorResponse — внутренний метод: собирает ErrorResponse и вызывает JSONResponse.
func (h *HTTPResponseHandler) errorResponse(
	statusCode int,
	err error,
	msg string,
) {
	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	h.JSONResponse(response, statusCode)

}
