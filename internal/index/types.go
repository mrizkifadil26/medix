package index

type IconIndexerConfig struct {
	Sources     []IconSource // List of source dirs
	OutputPath  string       // Where to save the final JSON
	ExcludeDirs []string     // Folder to exclude (e.g. "Collection")
}

type IconSource struct {
	Path   string
	Source string // Label, e.g., "personal", "downloaded"
}
