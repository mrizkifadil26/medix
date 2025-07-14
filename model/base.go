package model

import "time"

type MediaOutputLegacy[T any] struct {
	Type        string     `json:"type"`
	GeneratedAt time.Time  `json:"generated_at"`
	Data        []Group[T] `json:"data"`
}

type GroupLegacy[T any] struct {
	Name  string `json:"name"`
	Items []T    `json:"items"`
}
