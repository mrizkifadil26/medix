// Package utils provides general-purpose utility types and functions.
//
// This file defines custom flag.Value implementations used to enhance CLI argument parsing.
// Specifically, it includes:
//   - ArrayFlags: for handling comma-separated string lists (e.g., --include=a,b,c)
//   - KVFlags: for parsing key-value formatted strings (e.g., --label=prod:/path/to/data)
package utils

import (
	"fmt"
	"strings"
)

// ArrayFlags is a custom flag type that collects repeated string values from comma-separated input.
// Example: --include=alpha,beta,gamma
type ArrayFlags []string

// String returns the string representation of ArrayFlags.
func (a *ArrayFlags) String() string {
	return strings.Join(*a, ", ")
}

// Set parses a comma-separated string and appends non-empty values to the ArrayFlags list.
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

// KVFlags is a custom flag type that parses key-value pairs from input strings.
// Format: "key:path" (e.g., --label=data:/path/to/data)
// If no key is given (e.g., --label=/path/only), the path is used as both key and value.
type KVFlags struct {
	values map[string]string
	order  []string // Optional: preserve the insertion order of keys
}

// String returns the KVFlags in "key:path" comma-separated format.
func (k *KVFlags) String() string {
	if k.values == nil {
		return ""
	}
	var parts []string
	for _, key := range k.order {
		parts = append(parts, fmt.Sprintf("%s:%s", key, k.values[key]))
	}
	return strings.Join(parts, ",")
}

// Set parses a key:value string, handling both labeled and unlabeled paths.
// If only a single value is given, it is used as both key and value.
func (k *KVFlags) Set(value string) error {
	if k.values == nil {
		k.values = make(map[string]string)
	}
	// Try to split as label:path
	parts := strings.SplitN(value, ":", 2)

	if len(parts) == 2 {
		label := parts[0]
		path := parts[1]
		k.values[label] = path
		k.order = append(k.order, label)
	} else {
		// No label, fallback: use path as label
		path := parts[0]
		k.values[path] = path
		k.order = append(k.order, path)
	}
	return nil
}

// ToMap returns the parsed key-value pairs as a standard map.
func (k *KVFlags) ToMap() map[string]string {
	return k.values
}
