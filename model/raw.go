package model

import (
	"encoding/json"
	"fmt"
)

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

type IconMeta struct {
	ID       string `json:"id,omitempty"` // e.g. "sci-fi"
	Name     string `json:"name"`
	FullPath string `json:"full_path"`
	Size     int64  `json:"size"`
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
