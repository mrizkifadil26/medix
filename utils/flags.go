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
	*a = append(*a, value)
	return nil
}
