package model

// Top-level synced structure for media (movies / tvshows)
type SyncedOutput = MediaOutput[SyncedItem]
type SyncedGenre = Group[SyncedItem]

// One movie/tvshow folder entry (enriched version of RawItem)
type SyncedItem struct {
	Type   string            `json:"type"` // single or collection
	Name   string            `json:"name"`
	Path   string            `json:"path"`
	Status string            `json:"status"`
	Icon   *SyncedIconMeta   `json:"icon,omitempty"`   // Local .ico inside media folder
	Source *SyncedIconMeta   `json:"source,omitempty"` // Linked from index
	Items  []SyncedChildItem `json:"items,omitempty"`
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
	ID       string `json:"id,omitempty"` // For mapping
	Name     string `json:"name"`
	FullPath string `json:"full_path"`
	Size     int64  `json:"size"`
	Source   string `json:"source"` // "downloaded" or "personal"
	Type     string `json:"type"`   // "icon" or "collection"
}

// Synced version of Icon Index
type SyncedIconIndex = MediaOutput[SyncedIconEntry]
type SyncedIconGroup = Group[SyncedIconEntry]

// One icon (linked back to media if used)
type SyncedIconEntry struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	FullPath string           `json:"full_path"`
	Size     int64            `json:"size"`
	Source   string           `json:"source"`            // downloaded/personal
	Type     string           `json:"type"`              // icon / collection
	UsedBy   *UsedBy          `json:"used_by,omitempty"` // reverse link
	Variants []SyncedIconMeta `json:"variants,omitempty"`
}

type UsedBy struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	ContentType string `json:"content_type"` // "movies" or "tvshows"
}
