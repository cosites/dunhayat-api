package logger

import (
	"context"

	"dunhayat-api/pkg/logger"

	"go.uber.org/zap"
)

type Adapter struct {
	baseLogger logger.Interface
}

func NewAdapter(baseLogger logger.Interface) *Adapter {
	return &Adapter{
		baseLogger: baseLogger,
	}
}

func (a *Adapter) Info(msg string, fields ...zap.Field) {
	a.baseLogger.Info(msg, fields...)
}

func (a *Adapter) Error(msg string, fields ...zap.Field) {
	a.baseLogger.Error(msg, fields...)
}

func (a *Adapter) Warn(msg string, fields ...zap.Field) {
	a.baseLogger.Warn(msg, fields...)
}

func (a *Adapter) Debug(msg string, fields ...zap.Field) {
	a.baseLogger.Debug(msg, fields...)
}

func (a *Adapter) WithContext(ctx context.Context) *logger.Logger {
	return a.baseLogger.WithContext(ctx)
}

func (a *Adapter) WithAuthContext(
	ctx context.Context,
	userID string,
	action string,
) logger.Interface {
	loggerWithCtx := a.baseLogger.WithContext(ctx)

	return loggerWithCtx
}
