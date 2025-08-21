// utils/jsonutil.go
//
// Package utils provides general-purpose utility functions.
//
// This file contains helper functions for reading and writing JSON data
// to and from files, simplifying use of encoding/json with file operations.
package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// LoadJSON reads a JSON file from the given path and decodes it into the provided value.
// The value 'v' should be a pointer to the target struct, slice, or map.
//
// Example:
//
//	var config Config
//	err := utils.LoadJSON("config.json", &config)
func LoadJSON(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(v)
}

func LoadJSONOrdered(path string, om *OrderedMap[string, any]) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return om.UnmarshalJSON(data)
}

// WriteJSON writes the given data as pretty-formatted JSON to the specified file path.
// It creates parent directories if they don't exist.
//
// Example:
//
//	err := utils.WriteJSON("output.json", myData)
func WriteJSON(path string, data any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

/* func WriteJSONOrdered(path string, data *OrderedMap[string, any]) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
} */
