package utils

import (
	"fmt"
	"path/filepath"
)

// WriteReport writes any report struct to data/reports/*.report.json
func WriteReport(name string, v any) error {
	path := filepath.Join("data", "reports", fmt.Sprintf("%s.report.json", name))

	err := WriteJSON(path, v)
	if err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	fmt.Printf("📝 Report written: %s\n", path)
	return nil
}
