package util

import (
	"gopkg.in/yaml.v3"
	"os"
)

func YamlReader(filePath string, i interface{}) (err error) {
	byt, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(byt, i)
}
