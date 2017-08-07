package server

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

func (g *Logger) Load() error {
	temp, err := zap.NewProduction()
	if err != nil {
		return errors.Wrap(err, "cannot load logger")
	}
	g.logger = temp
	return nil
}

func (g *Logger) Err(m string, fields ...zapcore.Field) {
	g.logger.Error(m, fields...)
}

func (g *Logger) Info(m string, fields ...zapcore.Field) {
	g.logger.Info(m, fields...)
}

func (g *Logger) With(fields ...zapcore.Field) *Logger {
	return &Logger{logger: g.logger.With(fields...)}
}
