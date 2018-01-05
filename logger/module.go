package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Module struct {
	*zap.Logger
}

func (m *Module) Load() error {
	temp, err := zap.NewProduction(zap.AddCallerSkip(1))
	if err != nil {
		return fmt.Errorf("cannot load logger: %v", err)
	}
	m.Logger = temp
	return nil
}

func (m *Module) Err(s string, fields ...zapcore.Field) {
	m.Error(s, fields...)
}

func (m *Module) Info(s string, fields ...zapcore.Field) {
	m.Info(s, fields...)
}

func (m *Module) With(fields ...zapcore.Field) *Module {
	return &Module{Logger: m.With(fields...).Logger}
}
