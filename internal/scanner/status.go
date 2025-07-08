package scanner

import (
	"os"
	"path/filepath"
)

func ResolveStatus(path string) string {
	switch {
	case hasIcoFile(path) && fileExists(filepath.Join(path, "desktop.ini")):
		return "ok"
	case hasIcoFile(path):
		return "warn"
	default:
		return "missing"
	}
}

func hasIcoFile(dir string) bool {
	files, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".ico" {
			return true
		}
	}
	return false
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func isCollection(entries []os.DirEntry) bool {
	for _, e := range entries {
		if e.IsDir() {
			return true
		}
	}
	return false
}
