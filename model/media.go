package model

import "time"

type MediaOutput struct {
	Type           Type         `json:"type"`           // always "media"
	Version        string       `json:"version"`        // e.g. "1.0.0"
	GeneratedAt    time.Time    `json:"generatedAt"`    // generation timestamp
	Sources        []string     `json:"sources"`        // "movies" or "tv"
	TotalItems     int          `json:"totalItems"`     // number of entries
	GroupCount     int          `json:"groupCount"`     // number of groups (flat=genres, structured=top-level)
	ScanDurationMs int64        `json:"scanDurationMs"` // for performance
	Items          []MediaEntry `json:"items"`          // []MediaEntry or []MediaGroup
}

type MediaEntry struct {
	BaseEntry
	// Type   string       `json:"type"` // "single", "collection", "show", "season"
	Status string       `json:"status"`
	Icon   *IconRef     `json:"icon,omitempty"`
	Parent string       `json:"parent,omitempty"` // ✅ added
	Items  []MediaEntry `json:"items,omitempty"`
}

func (e MediaEntry) GetName() string { return e.Name }
