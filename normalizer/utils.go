package normalizer

import (
	"encoding/json"
	"fmt"
)

// toStringMap tries to convert map[string]any or map[string]string to map[string]string
func toStringMap(v any) (map[string]string, bool) {
	switch m := v.(type) {
	case map[string]string:
		// copy to avoid mutation of caller's map
		out := make(map[string]string, len(m))
		for k, val := range m {
			out[k] = val
		}

		return out, true
	case map[string]any:
		out := make(map[string]string, len(m))
		for k, val := range m {
			if s, ok := val.(string); ok {
				out[k] = s
			} else {
				// best-effort convert non-strings to string using fmt
				out[k] = fmt.Sprintf("%v", val)
			}
		}

		return out, true
	default:
		return nil, false
	}
}

func clone[K comparable, V any](src map[K]V) map[K]V {
	dst := make(map[K]V, len(src))
	for k, v := range src {
		dst[k] = v
	}

	return dst
}

func toAnySlice(strs []string) []any {
	out := make([]any, len(strs))
	for i, s := range strs {
		out[i] = s
	}

	return out
}

func PrettyPrintJSON(field map[string]any) string {
	bytes, err := json.MarshalIndent(field, "", "  ")
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(bytes)
}
