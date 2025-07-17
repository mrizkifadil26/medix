package model

type IconMeta struct {
	ID       string `json:"id,omitempty"` // e.g. "sci-fi"
	Name     string `json:"name"`
	FullPath string `json:"full_path"`
	Size     int64  `json:"size"`
}

type IconEntry struct {
	ID       string      `json:"id,omitempty"`
	Name     string      `json:"name"`
	FullPath string      `json:"full_path,omitempty"`
	Size     int64       `json:"size,omitempty"`
	Source   string      `json:"source,omitempty"`
	Type     string      `json:"type"` // "icon" or "collection"
	Items    []IconEntry `json:"items,omitempty"`
}
