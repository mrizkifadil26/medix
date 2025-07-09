package sync

import "time"

type UnusedIconReport struct {
	GeneratedAt time.Time         `json:"generated_at"`
	Total       int               `json:"total"`
	Icons       []UnusedIconEntry `json:"icons"`
}

type UnusedIconEntry struct {
	Name     string `json:"name"`
	FullPath string `json:"full_path"`
	Source   string `json:"source,omitempty"`
}
