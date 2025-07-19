package organizer

type Config struct {
	Mode       string `json:"mode"`       // "preview" or "apply"
	MediaInput string `json:"mediaInput"` // used in preview or fallback
	// ResultInput string `json:"resultInput"` // used in apply mode
	TargetDir   string   `json:"targetDir"`   // preview or apply result
	IconSources []string `json:"iconSources"` // list of icon directories
	// DryRun      bool     `json:"dryRun"`      // only for apply
}
