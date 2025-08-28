package logger

import (
	"context"

	"go.uber.org/zap"
)

type Interface interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	WithContext(ctx context.Context) *Logger
}

var _ Interface = (*Logger)(nil)
