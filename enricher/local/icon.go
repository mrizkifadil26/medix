package local

import (
	"fmt"
	"strings"

	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

type IconFilter struct{}

func (f IconFilter) Name() string {
	return "icon"
}

func (f IconFilter) Apply(
	item any,
	errs *[]error,
) {
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

	var icon *IconSource
	for _, node := range sources {
		nameVal, _ := jsonpath.Get(node, "name")
		extVal, _ := jsonpath.Get(node, "ext")
		pathVal, _ := jsonpath.Get(node, "path")
		sizeVal, _ := jsonpath.Get(node, "size")

		name, _ := nameVal.(string)
		ext, _ := extVal.(string)
		path, _ := pathVal.(string)
		size, _ := sizeVal.(float64)

		if name == "" || ext == "" {
			*errs = append(*errs, fmt.Errorf("invalid media source, missing name or ext"))
			continue
		}

		if !strings.EqualFold(ext, ".ico") {
			continue
		}

		icon = &IconSource{
			Name:      name,
			Extension: ext,
			Path:      path,
			Size:      int64(size),
		}
	}

	// Only set if we actually found one
	if icon != nil {
		_ = jsonpath.Set(item, "icon", icon)
	}
}
