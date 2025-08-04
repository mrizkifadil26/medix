package scannerV2

import (
	"os"
	"path/filepath"
	"strings"
)

func containsDir(entries []os.DirEntry) bool {
	for _, e := range entries {
		if e.IsDir() {
			return true
		}
	}
	return false
}

func pathDepth(base, path string) int {
	rel, _ := filepath.Rel(base, path)
	if rel == "." {
		return 0
	}

	return len(strings.Split(rel, string(filepath.Separator)))
}

func getLeafDepth(start string) (int, error) {
	level := 0
	current := start

	for {
		entries, err := os.ReadDir(current)
		if err != nil {
			return 0, err
		}
		found := false
		for _, e := range entries {
			if e.IsDir() {
				current = filepath.Join(current, e.Name())
				found = true
				break
			}
		}
		if !found {
			break
		}
		level++
	}

	return level, nil
}
