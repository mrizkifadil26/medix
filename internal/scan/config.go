// internal/scan/config.go
package scan

type ScanConfig struct {
	ContentType string   `json:"content_type"`       // "movies" or "tvshows"
	Sources     []string `json:"sources"`            // List of directories
	OutputPath  string   `json:"output_path"`        // Output file path
	Strategy    string   `json:"strategy,omitempty"` // (optional) for future use
}

type ScanConfigFile struct {
	Configs []ScanConfig `json:"configs"`
}
