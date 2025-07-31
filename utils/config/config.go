package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

// Parse loads and validates a config file into the provided target struct pointer
func Parse(configPath string, target interface{}) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	switch ext := filepath.Ext(configPath); ext {
	case ".json":
		err = json.Unmarshal(data, target)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, target)
	default:
		return fmt.Errorf("unsupported format: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	// Validate
	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(target); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}
