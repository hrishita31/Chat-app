package logger

import "go.uber.org/zap"

type ILogger interface {
	CreateNewLogger() *Logger
}

type Logger struct {
	L *zap.Logger
}

func NewLogger() *Logger {
	logger := zap.Must(zap.NewDevelopment())
	defer logger.Sync()
	logger.Info("hello from logger")

	return &Logger{
		L: logger,
	}
}
