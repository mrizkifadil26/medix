package model

type BaseEntry struct {
	Name   string    `json:"name"`
	Path   string    `json:"path"`
	Type   string    `json:"type"` // "single", "collection", "show", "season"
	Status string    `json:"status"`
	Icon   *IconMeta `json:"icon,omitempty"`
	Group  []string  `json:"group"`
	Parent string    `json:"parent,omitempty"` // ✅ added
}

type IconMeta struct {
	ID       string `json:"id,omitempty"` // e.g. "sci-fi"
	Name     string `json:"name"`
	FullPath string `json:"full_path"`
	Size     int64  `json:"size"`
}
