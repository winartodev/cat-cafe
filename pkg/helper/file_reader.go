package helper

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ReadYaml(path string, out interface{}) error {
	buf, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	if err = yaml.Unmarshal(buf, out); err != nil {
		return fmt.Errorf("unmarshal yaml: %w", err)
	}

	return nil
}
