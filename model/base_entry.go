package model

type BaseEntry struct {
	Name   string    `json:"name"`
	Path   string    `json:"path"`
	Type   string    `json:"type"` // "single", "collection", "show", "season"
	Status string    `json:"status"`
	Icon   *IconMeta `json:"icon,omitempty"`
	Group  string    `json:"group"`
	Parent string    `json:"parent,omitempty"` // ✅ added
}
