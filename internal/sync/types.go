package sync

import "time"

type UnusedIconReport struct {
	GeneratedAt time.Time                    `json:"generated_at"`
	Total       int                          `json:"total"`
	Groups      map[string][]UnusedIconEntry `json:"groups"`
}

type UnusedIconEntry struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Source string `json:"source,omitempty"`
}
