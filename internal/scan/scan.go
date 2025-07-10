package scan

const maxConcurrency = 8

func ScanAll[T any](cfg ScanConfig, strategy ScanStrategy[T]) T {
	return strategy.Scan(cfg.Sources)
}
