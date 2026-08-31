// Package core_logger предоставляет инструменты для логирования приложения.
package core_logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger представляет логгер приложения на основе zap.Logger.
type Logger struct {
	*zap.Logger

	file *os.File
}

type loggerContextKey struct{}

var (
	key = loggerContextKey{}
)

// ToContext кладёт логгер в контекст. Вызывается в middleware Logger,
// чтобы все последующие обработчики могли получить логгер с request_id.
func ToContext(ctx context.Context, log *Logger) context.Context {
	return context.WithValue(
		ctx,
		key,
		log,
	)
}

// FromContext извлекает логгер из контекста или вызывает panic, если логгер отсутствует или имеет некорректный тип.
func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value(key).(*Logger)
	if !ok {
		panic("no logger in context")
	}

	return log
}

// NewLogger создаёт и настраивает логгер на основе переданной конфигурации.
// Логи записываются одновременно в стандартный вывод и в файл.
// Возвращает ошибку, если не удалось настроить уровень логирования, создать директорию или открыть файл для записи.
func NewLogger(config LoggerConfig) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	// Создаем папку для логов
	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	// Время лога
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(
		config.Folder,
		fmt.Sprintf("%s.log", timestamp),
	)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	// создаем ядра - одно пишет в консоль, другое пишет в файл
	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),
	)

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{
		Logger: zapLogger,
		file:   logFile,
	}, nil
}

// Close закрывает файл, используемый логгером для записи логов.
// Если файл не удалось закрыть, ошибка выводится в стандартный вывод.
func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close application logger:", err)
	}
}

// With возвращает новый логгер с добавленными полями контекста для всех последующих записей.
// Переопределяем стандартный метод With у Zap logger, для возврата нашего логгера
func (l *Logger) With(field ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(field...),
		file:   l.file,
	}
}
