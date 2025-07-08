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
}

type RawChild struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Status string `json:"status"`
}
