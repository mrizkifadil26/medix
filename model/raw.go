package model

import (
	"encoding/json"
	"fmt"
)

type BaseEntry struct {
	Type   string    `json:"type"` // "single" or "collection"
	Name   string    `json:"name"`
	Path   string    `json:"path"`
	Status string    `json:"status"`
	Icon   *IconMeta `json:"icon,omitempty"` // ⬅️ new field
}

type IconMeta struct {
	ID       string `json:"id,omitempty"` // e.g. "sci-fi"
	Name     string `json:"name"`
	FullPath string `json:"full_path"`
	Size     int64  `json:"size"`
}

type MovieOutput = MediaOutput[MovieEntry]
type MovieGroup = Group[MovieEntry]
type MovieEntry struct {
	BaseEntry
	Items []MovieEntry `json:"items,omitempty"` // recursive
}

type TVShowOutput = MediaOutput[TVShowEntry]
type TVShowGroup = Group[TVShowEntry]
type TVShowEntry struct {
	BaseEntry
	Seasons []string `json:"seasons,omitempty"`
}

type RawOutput = MediaOutput[RawEntry]
type RawGenre = Group[RawEntry]

type RawEntry struct {
	Type   string         `json:"type"` // "single" or "collection"
	Name   string         `json:"name"`
	Path   string         `json:"path"`
	Status string         `json:"status"`
	Icon   *IconMeta      `json:"icon,omitempty"`  // ⬅️ new field
	Items  *RawEntryItems `json:"items,omitempty"` // Only for collections
}

type RawEntryItems struct {
	Entries []RawEntry
	Seasons []string
}

func (r *RawEntryItems) UnmarshalJSON(data []byte) error {
	var entries []RawEntry
	if err := json.Unmarshal(data, &entries); err == nil {
		r.Entries = entries
		return nil
	}

	var seasons []string
	if err := json.Unmarshal(data, &seasons); err == nil {
		r.Seasons = seasons
		return nil
	}

	return fmt.Errorf("items must be either []RawEntry or []string")
}
