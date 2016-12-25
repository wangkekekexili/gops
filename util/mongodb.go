package util

import (
	"fmt"
	"os"
)

type mongodbConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	URL      string `yaml:"url"`
	DB       string `yaml:"db"`
}

func BuildMongodbURI() (string, error) {
	config := mongodbConfig{}

	// Try building the URI from config file.
	if err := getConfig("mongodb.yaml", &config); err == nil {
		return fmt.Sprintf("mongodb://%s:%s@%s/%s", config.User, config.Password, config.URL, config.DB), nil
	}

	// Try building the URI from environment variables.
	return os.Getenv("mongodb_uri"), nil
}
