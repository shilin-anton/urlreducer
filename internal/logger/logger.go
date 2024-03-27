package logger

import (
	"go.uber.org/zap"
	"net/http"
)

// Log синглтон.
var Log *zap.Logger = zap.NewNop()

type (
	// Структура для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

// Initialize инициализирует логер.
func Initialize(level string) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	// устанавливаем синглтон
	Log = zl
	return nil
}

func RequestLogger(uri string, method string, duration string) {
	Log.Info("got incoming HTTP request",
		zap.String("URI", uri),
		zap.String("method", method),
		zap.String("duration", duration),
	)
}

// ResponseLogger — middleware-логер для HTTP-ответов.
func ResponseLogger(status string, size string) {
	Log.Info("HTTP response has been sent",
		zap.String("code", status),
		zap.String("size", size),
	)
}
