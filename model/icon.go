package model

import "time"

type IconIndex struct {
	Type           string      `json:"type"`           // Always "icon-index"
	Version        string      `json:"version"`        // Schema version, e.g. "1.0.0"
	GeneratedAt    time.Time   `json:"generatedAt"`    // When this index was created
	TotalItems     int         `json:"totalItems"`     // Total number of icon entries
	TotalSources   int         `json:"totalSources"`   // Count of distinct sources
	GroupCount     int         `json:"groupCount"`     // Count of top-level folder groups
	ScanDurationMs int64       `json:"scanDurationMs"` // e.g. "120ms", "2.3s"
	Items          []IconEntry `json:"data"`           // Flat list of icon entries
}

type IconEntry struct {
	ID       string   `json:"id,omitempty"`
	Name     string   `json:"name"`
	Size     int64    `json:"size,omitempty"`
	Source   string   `json:"source,omitempty"`
	Path     string   `json:"path,omitempty"`
	Group    []string `json:"group,omitempty"`
	Type     string   `json:"type"`               // "icon" or "collection"
	Variants []string `json:"variants,omitempty"` // List of full paths to variant .ico files
}
