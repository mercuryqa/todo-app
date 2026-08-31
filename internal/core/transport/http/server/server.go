// Package server содержит компоненты для запуска и настройки HTTP-сервера.
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	coreLogger "github.com/mercuryqa/todo-app/internal/core/logger"
	"github.com/mercuryqa/todo-app/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

// HTTPServer - Мультиплексор определяет через какие middleware должен пройти запрос
// и в какой обработчик его направить
type HTTPServer struct {
	mux    *http.ServeMux
	config Config
	log    *coreLogger.Logger

	middleware []core_http_middleware.Middleware
}

// NewHTTPServer - конструктор для http сервера
func NewHTTPServer(
	config Config,
	log *coreLogger.Logger,
	middleware ...core_http_middleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		config:     config,
		log:        log,
		middleware: middleware,
	}
}

func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router.WithMiddleware()))
	}
}

// Run - запускает сервер
func (s *HTTPServer) Run(ctx context.Context) error {
	// навешиваем на мультиплексор цепь middleware
	mux := core_http_middleware.ChainMiddleware(s.mux, s.middleware...)

	server := http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("start HTTP server", zap.String("addr", s.config.Addr))

		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}
		s.log.Warn("HTTP server stopped")
	}
	return nil
}
