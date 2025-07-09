package index

type IconIndexerUnifiedConfig struct {
	Outputs     map[string]OutputTarget `json:"outputs"`      // key: category (e.g. movies, tvshows)
	ExcludeDirs []string                `json:"exclude_dirs"` // shared across groups
}

type OutputTarget struct {
	OutputPath string       `json:"output_path"`
	Sources    []IconSource `json:"sources"`
}

type IconSource struct {
	Path   string `json:"path"`
	Source string `json:"source"`
}

type IconIndexerConfig struct {
	Sources     []IconSource `json:"sources"`      // Input folders for a single content type (e.g. movies)
	OutputPath  string       `json:"output_path"`  // Where to save the generated JSON
	ExcludeDirs []string     `json:"exclude_dirs"` // Optional: skip these subfolders during scan
}
