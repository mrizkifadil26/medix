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
	Type     string    `json:"type"` // "single" or "collection"
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Status   string    `json:"status"`
	Children any       `json:"items,omitempty"` // Only for collections
	Icon     *IconMeta `json:"icon,omitempty"`  // ⬅️ new field
}

type RawChild struct {
	Type   string    `json:"type"` // "single" or "collection"
	Name   string    `json:"name"`
	Path   string    `json:"path"`
	Status string    `json:"status"`
	Icon   *IconMeta `json:"icon,omitempty"` // ⬅️ new field
}

type IconMeta struct {
	Name     string `json:"name"`
	FullPath string `json:"full_path"`
	Size     int64  `json:"size"`
}

type IconIndex struct {
	Type        string      `json:"type"` // "genre"
	GeneratedAt time.Time   `json:"generated_at"`
	Groups      []IconGroup `json:"groups"`
}

type IconGroup struct {
	ID    string      `json:"id,omitempty"` // e.g. "sci-fi"
	Name  string      `json:"name"`         // genre name like "Sci-Fi"
	Items []IconEntry `json:"items"`
}

type IconEntry struct {
	ID        string      `json:"id,omitempty"`
	Name      string      `json:"name"`
	Extension string      `json:"extension,omitempty"`
	Size      int64       `json:"size,omitempty"`
	Source    string      `json:"source,omitempty"`
	FullPath  string      `json:"full_path,omitempty"`
	Type      string      `json:"type"` // "icon" or "collection"
	Items     []IconEntry `json:"items,omitempty"`
}
