package model

import "time"

type MediaOutput struct {
	Type           string         `json:"type"`           // "flat" or "structured"
	Version        string         `json:"version"`        // e.g. "1.0.0"
	GeneratedAt    time.Time      `json:"generatedAt"`    // generation timestamp
	Source         string         `json:"source"`         // "movies" or "tvshows"
	TotalItems     int            `json:"totalItems"`     // number of entries
	GroupCount     int            `json:"groupCount"`     // number of groups (flat=genres, structured=top-level)
	ScanDurationMs int64          `json:"scanDurationMs"` // for performance
	Items          []MediaEntry   `json:"items"`          // []MediaEntry or []MediaGroup
}
