package scannerV2

const (
	DefaultOutputPath = "output.json"
	DefaultDepth      = 1
	DefaultMode       = "mixed"
)

func DefaultConfig() Config {
	return Config{
		Options: &ScanOptions{
			Mode:  DefaultMode,
			Depth: DefaultDepth,
		},
		Output: &OutputOptions{
			Format:          ptr("json"),
			OutputPath:      nil,
			IncludeErrors:   false,
			IncludeWarnings: false,
			IncludeStats:    false,
		},
	}
}

func ptr[T any](v T) *T { return &v }
