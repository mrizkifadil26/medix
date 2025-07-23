package model

import "time"

type IconIndex struct {
	Type           Type        `json:"type"`           // Always "icon"
	Version        string      `json:"version"`        // Schema version, e.g. "1.0.0"
	GeneratedAt    time.Time   `json:"generatedAt"`    // When this index was created
	Sources        []string    `json:"sources"`        // "movies" or "tv"
	TotalItems     int         `json:"totalItems"`     // Total number of icon entries
	GroupCount     int         `json:"groupCount"`     // Count of top-level folder groups
	ScanDurationMs int64       `json:"scanDurationMs"` // e.g. "120ms", "2.3s"
	Items          []IconEntry `json:"data"`           // Flat list of icon entries
}

type IconEntry struct {
	BaseEntry
	Slug     string   `json:"id"`
	Size     int64    `json:"size"`
	Variants []string `json:"variants,omitempty"` // List of full paths to variant .ico files
}
