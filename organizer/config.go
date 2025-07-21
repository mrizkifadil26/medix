package organizer

type Config struct {
	MediaInput  string       `json:"mediaInput"` // used in preview or fallback
	TargetDir   string       `json:"targetDir"`  // preview or apply result
	ReportPath  string       `json:"reportPath"`
	IconSources []IconSource `json:"iconSources"` // list of icon directories
}

type IconSource struct {
	Path   string `json:"path"`
	Source string `json:"source"`
}
