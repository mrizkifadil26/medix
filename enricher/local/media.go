package local

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

type MediaFilter struct{}

func (f MediaFilter) Name() string { return "media" }

func (f MediaFilter) Apply(item any, errs *[]error) {
	var sources []any

	// Determine if item is a directory
	typeVal, _ := jsonpath.Get(item, "type")
	if t, ok := typeVal.(string); ok && t == "directory" {
		childrenNode, err := jsonpath.Get(item, "children")
		if err != nil {
			*errs = append(*errs, fmt.Errorf("directory has no children"))
			return
		}

		if childrenArr, ok := childrenNode.([]any); ok {
			sources = childrenArr
		} else {
			*errs = append(*errs, fmt.Errorf("children is not an array"))
			return
		}
	} else {
		// Single file or no type
		sources = []any{item}
	}

	var mainMedia *MediaSource
	var extras []*MediaSource

	for _, node := range sources {
		nameVal, _ := jsonpath.Get(node, "name")
		extVal, _ := jsonpath.Get(node, "ext")
		pathVal, _ := jsonpath.Get(node, "path")
		sizeVal, _ := jsonpath.Get(node, "size")

		name, _ := nameVal.(string)
		ext, _ := extVal.(string)
		path, _ := pathVal.(string)
		size, _ := sizeVal.(float64)

		// Fallback: extract extension from name if ext is empty
		if ext == "" && name != "" {
			ext = strings.ToLower(filepath.Ext(name))
		}

		if name == "" || ext == "" {
			*errs = append(*errs, fmt.Errorf("invalid media source, missing name or ext"))
			continue
		}

		ms := &MediaSource{
			Name:      name,
			Extension: ext,
			Path:      path,
			Size:      int64(size),
		}

		if !strings.EqualFold(ext, ".mkv") {
			continue
		}

		// Decide main vs extras
		if mainMedia == nil {
			mainMedia = ms
		} else {
			lname := strings.ToUpper(name)
			if strings.Contains(lname, "IMAX") || strings.Contains(lname, "OPEN MATTE") {
				extras = append(extras, ms)
			} else {
				extras = append(extras, ms)
			}
		}
	}

	if mainMedia != nil {
		_ = jsonpath.Set(item, "media", mainMedia)
	}

	if len(extras) > 0 {
		_ = jsonpath.Set(item, "extras", extras)
	}
}
