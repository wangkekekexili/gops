package server

import (
	"os"
	"strings"
	"reflect"

	"github.com/pkg/errors"
)

type Config struct {
	MysqlDSN   string `env:"MYSQL_DSN,required"`
	SentryDSN  string `env:"SENTRY_DSN,required"`
	TillURL    string `env:"TILL_URL"`
	TillTarget string `env:"TILL_TARGET"`
}

func (c *Config) Load() error {
	configValue := reflect.ValueOf(c).Elem()
	configType := reflect.TypeOf(c).Elem()
	for i := 0; i != configValue.NumField(); i++ {
		configFieldValue := configValue.Field(i)
		configFieldType := configType.Field(i)
		envInfo, ok := configFieldType.Tag.Lookup("env")
		if !ok || envInfo == "" {
			continue
		}
		parts := strings.SplitN(envInfo, ",", 2)

		envKey := parts[0]
		var required bool
		if len(parts) == 2 {
			if parts[1] == "required" {
				required = true
			}
		}

		envValue := strings.TrimSpace(os.Getenv(envKey))
		if envValue == "" && required {
			return errors.Errorf("%v is required", envKey)
		}

		configFieldValue.SetString(envValue)
	}
	return nil
}
