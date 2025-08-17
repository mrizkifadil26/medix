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

	"github.com/iancoleman/orderedmap"
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

func LoadJSONOrdered(path string, v *orderedmap.OrderedMap) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// var tmp map[string]any
	// if err := json.Unmarshal(data, &tmp); err != nil {
	// 	return err
	// }

	// convertIntoOrderedMap(tmp, v)
	// return nil
	return v.UnmarshalJSON(data)
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

// convertIntoOrderedMap fills dst with the contents of src, recursively converting maps
func convertIntoOrderedMap(src map[string]any, dst *orderedmap.OrderedMap) {
	for k, val := range src {
		dst.Set(k, convertValue(val))
	}
}

// convertValue converts maps to OrderedMap recursively, arrays to []any containing OrderedMaps if needed
func convertValue(val any) any {
	switch v := val.(type) {
	case map[string]any:
		omap := orderedmap.New()
		convertIntoOrderedMap(v, omap)
		return omap
	case []any:
		for i := range v {
			v[i] = convertValue(v[i])
		}
		return v
	default:
		return v
	}
}
