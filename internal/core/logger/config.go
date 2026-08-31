// Package core_logger предоставляет инструменты для логирования приложения.
package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// LoggerConfig содержит конфигурацию логгера.
type LoggerConfig struct {
	Level  string `envconfig:"LEVEL" default:"DEBUG"`
	Folder string `envconfig:"FOLDER" required:"true"`
}

// NewLoggerConfig возвращает конфигурацию логгера, полученную из переменных окружения,
// или ошибку, если не удалось обработать конфигурацию.
func NewLoggerConfig() (LoggerConfig, error) {
	var config LoggerConfig

	if err := envconfig.Process("LOGGER", &config); err != nil {
		return LoggerConfig{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

// NewCLoggerConfigMust возвращает конфигурацию логгера или вызывает panic, если не удалось её создать.
func NewCLoggerConfigMust() LoggerConfig {
	config, err := NewLoggerConfig()
	if err != nil {
		fmt.Errorf("new logger congig: %w", err)
		panic(err)
	}
	return config
}
