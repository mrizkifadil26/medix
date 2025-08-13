package scanner

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
			Format:          "json",
			OutputPath:      DefaultOutputPath,
			IncludeErrors:   false,
			IncludeWarnings: false,
			IncludeStats:    false,
		},
	}
}
