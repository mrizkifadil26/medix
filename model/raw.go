package model

type RawOutput = MediaOutput[RawEntry]
type RawGenre = Group[RawEntry]

type RawEntry struct {
	Type   string     `json:"type"` // "single" or "collection"
	Name   string     `json:"name"`
	Path   string     `json:"path"`
	Status string     `json:"status"`
	Icon   *IconMeta  `json:"icon,omitempty"`  // ⬅️ new field
	Items  []RawEntry `json:"items,omitempty"` // Only for collections
}

type IconMeta struct {
	ID       string `json:"id,omitempty"` // e.g. "sci-fi"
	Name     string `json:"name"`
	FullPath string `json:"full_path"`
	Size     int64  `json:"size"`
}
