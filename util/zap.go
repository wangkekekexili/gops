package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	// Init logger.
	logger, _ = zap.NewProduction()
}

func LogError(message string, fields ...zapcore.Field) {
	logger.Error(message, fields...)
}

func LogInfo(message string, fields ...zapcore.Field) {
	logger.Info(message, fields...)
}
