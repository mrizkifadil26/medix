package normalizer

import (
	"path/filepath"
	"strings"
)

func StripExtension(s string) string {
	return strings.TrimSuffix(s, filepath.Ext(s))
}
