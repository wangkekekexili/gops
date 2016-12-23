package util

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func getConfig(config string, target interface{}) error {
	bytes, err := ioutil.ReadFile(fmt.Sprintf("config%c%s", filepath.Separator, config))
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, target)
	return err
}
