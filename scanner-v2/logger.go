package scannerV2

const (
	LevelError = 1
	LevelInfo  = 2
	LevelDebug = 3
	LevelTrace = 4
)

func getLogLevel(opts ScanOptions) int {
	if opts.Trace {
		return LevelTrace
	}

	if opts.Debug {
		return LevelDebug
	}

	if opts.Verbose {
		return LevelInfo
	}

	return LevelError
}
