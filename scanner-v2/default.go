package scannerV2

const (
	DefaultOutputPath = "output.json"
	DefaultDepth      = 1
	DefaultMode       = "mixed"
)

func DefaultConfig() Config {
	return Config{
		Verbose: false,
		Options: ScanOptions{
			Mode:  DefaultMode,
			Depth: DefaultDepth,
		},
		Output: OutputOptions{
			Format:          "json",
			OutputPath:      nil,
			IncludeErrors:   false,
			IncludeWarnings: false,
			IncludeStats:    false,
		},
	}
}
