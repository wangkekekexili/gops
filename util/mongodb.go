package util

import "fmt"

type mongodbConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	URL      string `yaml:"url"`
	DB       string `yaml:"db"`
}

func BuildMongodbURI() (string, error) {
	config := mongodbConfig{}
	if err := getConfig("mongodb.yaml", &config); err != nil {
		return "", err
	}
	return fmt.Sprintf("mongodb://%s:%s@%s/%s", config.User, config.Password, config.URL, config.DB), nil
}
