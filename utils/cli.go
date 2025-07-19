package utils

import (
	"strings"
)

// arrayFlags is a custom flag type for collecting repeated string values.
type ArrayFlags []string

func (a *ArrayFlags) String() string {
	return strings.Join(*a, ", ")
}

func (a *ArrayFlags) Set(value string) error {
	parts := strings.Split(value, ",")
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			*a = append(*a, trimmed)
		}
	}

	return nil
}
