package model

type IconMeta struct {
	ID       string `json:"id,omitempty"` // e.g. "sci-fi"
	Name     string `json:"name"`
	FullPath string `json:"full_path"`
	Size     int64  `json:"size"`
}
