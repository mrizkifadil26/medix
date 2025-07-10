package organize

type SourceDir struct {
	Path   string `json:"path"`
	Source string `json:"source"`
}

type OrganizeConfig struct {
	RawMetadataPath string      `json:"raw_metadata_path"`
	Sources         []SourceDir `json:"sources"`
	OutputBase      string      `json:"output_base"`
	ExcludeDirs     []string    `json:"exclude_dirs"`
	PlanPath        string      `json:"plan_path"`
}
