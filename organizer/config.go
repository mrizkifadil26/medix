package organizer

type Config struct {
	MediaInput string `json:"mediaInput"` // used in preview or fallback
	// ResultInput string `json:"resultInput"` // used in apply mode
	TargetDir   string       `json:"targetDir"`   // preview or apply result
	IconSources []IconSource `json:"iconSources"` // list of icon directories
	// DryRun      bool     `json:"dryRun"`      // only for apply
}

type IconSource struct {
	Path   string `json:"path"`
	Source string `json:"source"`
}
