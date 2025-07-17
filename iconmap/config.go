package iconmap

import (
	"errors"
	"fmt"
)

type SourceConfig struct {
	Name string `json:"name"` // e.g. "downloaded"
	Type string `json:"type"` // "movies" | "tv"
	Path string `json:"path"`
}

type Config struct {
	Type        string         `json:"type"` // "movies" or "tv"
	Sources     []SourceConfig `json:"sources"`
	ExcludeDirs []string       `json:"excludeDirs"`
	OutputPath  string         `json:"output"`
}

// Validate ensures the config is sane.
func (c Config) Validate() error {
	if c.Type != "movies" && c.Type != "tv" {
		return fmt.Errorf("invalid media type: %q (must be 'movies' or 'tv')", c.Type)
	}

	if len(c.Sources) == 0 {
		return errors.New("at least one source must be provided")
	}

	for _, s := range c.Sources {
		if s.Path == "" || s.Name == "" {
			return errors.New("each source must have both path and name")
		}
	}

	if c.OutputPath == "" {
		return errors.New("output path must be specified")
	}

	return nil
}
