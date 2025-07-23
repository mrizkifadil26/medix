package model

type Type string

const (
	TypeMedia ContentType = "media"
	TypeIcon  ContentType = "icon"
)

type ContentType string

const (
	TypeMovies ContentType = "movies"
	TypeTV     ContentType = "tv"
)

type BaseEntry struct {
	Name        string      `json:"name"`
	Path        string      `json:"path"`
	Type        string      `json:"type"`        // "media" or "icon"
	ContentType ContentType `json:"contentType"` // "movies" or "tv"
	Group       []string    `json:"group"`
	Source      string      `json:"source"`
}

type IconRef struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Size     int64  `json:"size"`
	FullPath string `json:"fullPath"`
}
