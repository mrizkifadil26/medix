package model

type IconIndex = MediaOutput[IconEntry]
type IconGroup = Group[IconEntry]

type IconEntry struct {
	ID       string      `json:"id,omitempty"`
	Name     string      `json:"name"`
	FullPath string      `json:"full_path,omitempty"`
	Size     int64       `json:"size,omitempty"`
	Source   string      `json:"source,omitempty"`
	Type     string      `json:"type"` // "icon" or "collection"
	Items    []IconEntry `json:"items,omitempty"`
}
