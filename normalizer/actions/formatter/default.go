package formatter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mrizkifadil26/medix/utils/jsonpath"
)

func DefaultFormatter(input any, template string) (string, error) {
	result := template

	keys := extractPlaceholders(template)
	for _, key := range keys {
		val, err := jsonpath.Get(input, key)
		if err != nil || val == nil || !isPrimitive(val) {
			// Use placeholder for unknown values
			val = "[unknown]"
		}

		strVal := fmt.Sprintf("%v", val)
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, strVal)
	}

	return result, nil
}

func extractPlaceholders(template string) []string {
	// Match anything inside {{…}}
	re := regexp.MustCompile(`{{\s*([^{}]+)\s*}}`)
	matches := re.FindAllStringSubmatch(template, -1)

	keys := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) > 1 {
			keys = append(keys, m[1]) // m[1] is the content inside {{}}
		}
	}

	return keys
}

func isPrimitive(val any) bool {
	switch val.(type) {
	case string, bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return true
	default:
		return false
	}
}

func init() {
	GetRegistry().
		Register("default", DefaultFormatter)
}
