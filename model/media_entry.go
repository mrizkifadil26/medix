package model

type MediaEntry struct {
	BaseEntry
	Items []MediaEntry `json:"items,omitempty"`
}

func (e MediaEntry) GetName() string { return e.Name }
