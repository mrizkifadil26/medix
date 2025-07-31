package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mrizkifadil26/medix/utils/config"
	"github.com/stretchr/testify/assert"
)

type SampleConfig struct {
	Name   string `json:"name" yaml:"name" validate:"required"`
	Age    int    `json:"age" yaml:"age" validate:"gte=1,lte=120"`
	Active bool   `json:"active" yaml:"active"`
}

func writeTempFile(t *testing.T, content string, ext string) string {
	t.Helper()
	tmp := t.TempDir()
	path := filepath.Join(tmp, "config"+ext)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return path
}

func TestParseValidJSONConfig(t *testing.T) {
	content := `{
		"name": "Alice",
		"age": 30,
		"active": true
	}`
	path := writeTempFile(t, content, ".json")

	var cfg SampleConfig
	err := config.Parse(path, &cfg)

	assert.NoError(t, err)
	assert.Equal(t, "Alice", cfg.Name)
	assert.Equal(t, 30, cfg.Age)
	assert.Equal(t, true, cfg.Active)
}

func TestParseValidYAMLConfig(t *testing.T) {
	content := `
name: Bob
age: 25
active: false
`
	path := writeTempFile(t, content, ".yaml")

	var cfg SampleConfig
	err := config.Parse(path, &cfg)

	assert.NoError(t, err)
	assert.Equal(t, "Bob", cfg.Name)
	assert.Equal(t, 25, cfg.Age)
	assert.Equal(t, false, cfg.Active)
}

func TestMissingRequiredField(t *testing.T) {
	content := `
age: 25
active: true
`
	path := writeTempFile(t, content, ".yaml")

	var cfg SampleConfig
	err := config.Parse(path, &cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation error")
}

func TestInvalidFormatExtension(t *testing.T) {
	content := `name: John`
	path := writeTempFile(t, content, ".txt")

	var cfg SampleConfig
	err := config.Parse(path, &cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported format")
}

func TestInvalidYAMLFormat(t *testing.T) {
	content := `
name: Jane
age: not_a_number
`
	path := writeTempFile(t, content, ".yaml")

	var cfg SampleConfig
	err := config.Parse(path, &cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parse config")
}
