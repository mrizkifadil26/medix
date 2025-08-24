package transformer

import (
	"path/filepath"
	"strings"
)

func StripExtension(s string) (string, error) {
	return strings.TrimSuffix(s, filepath.Ext(s)), nil
}

func init() {
	GetRegistry().
		Register("stripExtension", StripExtension)
}
