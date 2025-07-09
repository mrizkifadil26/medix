package model

import "time"

type MediaOutput[T any] struct {
	Type        string     `json:"type"`
	GeneratedAt time.Time  `json:"generated_at"`
	Data        []Group[T] `json:"data"`
}

type Group[T any] struct {
	Name  string `json:"name"`
	Items []T    `json:"items"`
}
