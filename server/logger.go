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

func (g *Logger) Err(m string, fields map[string]interface{}) {
	g.logger.Error(m, g.mapToFields(fields)...)
}

func (g *Logger) Info(m string, fields map[string]interface{}) {
	g.logger.Info(m, g.mapToFields(fields)...)
}

func (g *Logger) InfoZap(m string, fields ...zapcore.Field) {
	g.logger.Info(m, fields...)
}

func (Logger) mapToFields(m map[string]interface{}) []zapcore.Field {
	var zapFields []zapcore.Field
	for key, valueI := range m {
		f := zapcore.Field{Key: key}
		switch value := valueI.(type) {
		case int64:
			f.Type = zapcore.Int64Type
			f.Integer = value
		case string:
			f.Type = zapcore.StringType
			f.String = value
		default:
			f.Interface = value
		}
		zapFields = append(zapFields, f)
	}
	return zapFields
}
