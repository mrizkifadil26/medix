package model

import "time"

// Top-level synced structure for media (movies / tvshows)
type SyncedOutput struct {
	Type        string        `json:"type"` // "movies" or "tvshows"
	GeneratedAt time.Time     `json:"generated_at"`
	Data        []SyncedGenre `json:"data"`
}

// Genre grouping
type SyncedGenre struct {
	Genre string       `json:"genre"`
	Items []SyncedItem `json:"items"`
}

// One movie/tvshow folder entry (enriched version of RawItem)
type SyncedItem struct {
	Type     string            `json:"type"` // single or collection
	Name     string            `json:"name"`
	Path     string            `json:"path"`
	Status   string            `json:"status"`
	Icon     *SyncedIconMeta   `json:"icon,omitempty"`   // Local .ico inside media folder
	Source   *SyncedIconMeta   `json:"source,omitempty"` // Linked from index
	Children []SyncedChildItem `json:"children,omitempty"`
}

// A season or subfolder
type SyncedChildItem struct {
	Type   string `json:"type,omitempty"` // Optional if you want to tag child as single
	Name   string `json:"name"`
	Path   string `json:"path"`
	Status string `json:"status"`
}

// Icon metadata (used for both local icon and source)
type SyncedIconMeta struct {
	ID        string `json:"id,omitempty"` // For mapping
	Name      string `json:"name"`
	Extension string `json:"extension,omitempty"`
	Size      int64  `json:"size"`
	Source    string `json:"source"` // "downloaded" or "personal"
	FullPath  string `json:"full_path"`
	Type      string `json:"type"` // "icon" or "collection"
}

// Synced version of Icon Index
type SyncedIconIndex struct {
	Type        string            `json:"type"` // "genre"
	GeneratedAt time.Time         `json:"generated_at"`
	Groups      []SyncedIconGroup `json:"groups"`
}

type SyncedIconGroup struct {
	ID    string            `json:"id"`
	Name  string            `json:"name"`
	Items []SyncedIconEntry `json:"items"`
}

// One icon (linked back to media if used)
type SyncedIconEntry struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Extension string  `json:"extension"`
	Size      int64   `json:"size"`
	Source    string  `json:"source"` // downloaded/personal
	FullPath  string  `json:"full_path"`
	Type      string  `json:"type"`              // icon / collection
	UsedBy    *UsedBy `json:"used_by,omitempty"` // reverse link
}

type UsedBy struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	ContentType string `json:"content_type"` // "movies" or "tvshows"
}
