package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func LoadConfig[T any](path string) (T, error) {
	var out T

	f, err := os.Open(path)
	if err != nil {
		return out, fmt.Errorf("open config file: %w", err)
	}
	defer f.Close()

	switch ext := strings.ToLower(filepath.Ext(path)); ext {
	case ".json":
		if err := json.NewDecoder(f).Decode(&out); err != nil {
			return out, fmt.Errorf("decode JSON: %w", err)
		}
	case ".yaml", ".yml":
		if err := yaml.NewDecoder(f).Decode(&out); err != nil {
			return out, fmt.Errorf("decode YAML: %w", err)
		}
	default:
		return out, fmt.Errorf("unsupported config file extension: %s", ext)
	}

	return out, nil
}
