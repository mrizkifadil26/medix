package scanner_test

import (
	"testing"

	"github.com/mrizkifadil26/medix/scanner"
)

func TestApplyDefaults(t *testing.T) {
	cfg := scanner.Config{}
	cfg.ApplyDefaults()

	if cfg.Options.Mode != "files" {
		t.Errorf("Expected Mode to default to 'files', got '%s'", cfg.Options.Mode)
	}

	if cfg.Options.Depth != 1 {
		t.Errorf("Expected Depth to default to 1, got %d", cfg.Options.Depth)
	}
}
