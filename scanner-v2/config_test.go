package scannerV2_test

import (
	"testing"

	scannerV2 "github.com/mrizkifadil26/medix/scanner-v2"
)

func TestApplyDefaults(t *testing.T) {
	cfg := scannerV2.Config{}
	cfg.ApplyDefaults()

	if cfg.Options.Mode != "files" {
		t.Errorf("Expected Mode to default to 'files', got '%s'", cfg.Options.Mode)
	}

	if cfg.Options.Depth != 1 {
		t.Errorf("Expected Depth to default to 1, got %d", cfg.Options.Depth)
	}
}
