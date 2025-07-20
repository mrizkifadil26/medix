package organizer

import (
	"time"

	"github.com/mrizkifadil26/medix/model"
)

type OrganizeResult struct {
	Type        string       `json:"type"`         // "tv" or "movies"
	GeneratedAt time.Time    `json:"generated_at"` // ISO8601 timestamp
	Changes     []ChangeItem `json:"changes"`      // List of file actions
}

type ChangeItem struct {
	Action string          `json:"action"` // "move", "copy", etc.
	Source string          `json:"source"` // Full path to original icon
	Target string          `json:"target"` // Target path to move/copy to
	Reason string          `json:"reason"` // Explanation of match
	Item   model.BaseEntry `json:"item"`   // Metadata about the media
}

type MediaEntry struct {
	ID    string `json:"id"`    // Unique show/movie ID
	Name  string `json:"name"`  // Title
	Group string `json:"group"` // Genre or classification
	Type  string `json:"type"`  // "show" or
}
