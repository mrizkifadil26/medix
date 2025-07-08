package model

import "time"

type RawOutput struct {
	Type        string       `json:"type"`         // "movies" or "tvshows"
	GeneratedAt time.Time    `json:"generated_at"` // Timestamp
	Data        []GenreBlock `json:"data"`         // Grouped by genre
}

type GenreBlock struct {
	Genre string    `json:"genre"`
	Items []RawItem `json:"items"`
}

type RawItem struct {
	Type     string     `json:"type"` // "single" or "collection"
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	Status   string     `json:"status"`
	Children []RawChild `json:"items,omitempty"` // Only for collections
	Icon     *IconMeta  `json:"icon,omitempty"`  // ⬅️ new field
}

type RawChild struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Status string `json:"status"`
}

type IconMeta struct {
	Name     string `json:"name"`
	FullPath string `json:"full_path"`
	Size     int64  `json:"size"`
}
