package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Module struct {
	logger *zap.Logger
}

func (m *Module) Load() error {
	temp, err := zap.NewProduction(zap.AddCallerSkip(1))
	if err != nil {
		return fmt.Errorf("cannot load logger: %v", err)
	}
	m.logger = temp
	return nil
}

func (m *Module) Err(s string, fields ...zapcore.Field) {
	m.logger.Error(s, fields...)
}

func (m *Module) Info(s string, fields ...zapcore.Field) {
	m.logger.Info(s, fields...)
}

func (m *Module) With(fields ...zapcore.Field) *Module {
	return &Module{logger: m.logger.With(fields...)}
}
