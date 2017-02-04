package util

import "github.com/uber-go/zap"

var logger zap.Logger

func init() {
	// Init logger.
	logger = zap.New(zap.NewJSONEncoder(zap.NoTime()))
}

func LogError(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

func LogInfo(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}
