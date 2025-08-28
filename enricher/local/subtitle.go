package local

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

type SubtitlesFilter struct{}

func (f SubtitlesFilter) Name() string {
	return "subtitle"
}

func (f SubtitlesFilter) Apply(
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

	subtitles := make(Subtitles)
	for _, child := range sources {
		nameVal, _ := jsonpath.Get(child, "name")
		extVal, _ := jsonpath.Get(child, "ext")
		pathVal, _ := jsonpath.Get(child, "path")
		// sizeVal, _ := jsonpath.Get(child, "size")

		name, _ := nameVal.(string)
		ext, _ := extVal.(string)
		path, _ := pathVal.(string)
		// size, _ := sizeVal.(float64)

		// Fallback: extract extension from name if ext is empty
		if ext == "" && name != "" {
			ext = strings.ToLower(filepath.Ext(name))
		}

		if name == "" || ext == "" {
			*errs = append(*errs, fmt.Errorf("invalid subtitle source, missing name or ext"))
			continue
		}

		// only accept .srt
		if ext != ".srt" {
			continue
		}

		// detect language
		lang := detectLangWithFilter(name, []string{"Pahe.in"})
		if lang == "" {
			lang = "id" // default
		}

		subtitles[lang] = Subtitle{
			Name: name,
			Path: path,
			Ext:  ext,
		}
	}

	if len(subtitles) > 0 {
		_ = jsonpath.Set(item, "subtitles", subtitles)
	}
}

func detectLangWithFilter(filename string, excludeFilters []string) string {
	lower := strings.ToLower(filename)

	// Strip .srt extension
	lower = strings.TrimSuffix(lower, ".srt")

	// Remove any exclude substrings
	for _, excl := range excludeFilters {
		lower = strings.ReplaceAll(lower, strings.ToLower(excl), "")
	}

	// Clean leftover separators
	lower = strings.Trim(lower, ".-_ ")

	// Look for language code after last dot or dash
	for _, sep := range []string{".", "-"} {
		if idx := strings.LastIndex(lower, sep); idx != -1 {
			part := lower[idx+1:]
			if len(part) == 2 { // simple ISO code check
				return part
			}
		}
	}

	// If nothing found, fallback
	return "id"
}
