// Package server содержит компоненты для запуска и настройки HTTP-сервера.
package server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config содержит конфигурацию HTTP-сервера.
type Config struct {
	Addr            string        `envconfig:"ADDR" required:"true"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"30s"`
}

// NewConfig возвращает конфигурацию HTTP-сервера, полученную из переменных окружения,
// или ошибку, если не удалось обработать конфигурацию.
func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("HTTP", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

// NewConfigMust возвращает конфигурацию HTTP-сервера или вызывает panic, если не удалось её создать.
func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get HTTP server config: %w", err)
		panic(err)
	}

	return config
}
