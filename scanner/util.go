package scanner

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/mrizkifadil26/medix/model"
	util "github.com/mrizkifadil26/medix/utils"
)

var globalConcurrency = 1

func SetConcurrency(n int) {
	if n > 0 {
		globalConcurrency = n
	}
}

func getConcurrency() int {
	return globalConcurrency
}

func findIcon(dir string, entries []os.DirEntry) *model.IconMeta {
	for _, f := range entries {
		if f.IsDir() || filepath.Ext(f.Name()) != ".ico" {
			continue
		}

		info, err := os.Stat(filepath.Join(dir, f.Name()))
		if err != nil {
			continue
		}

		return &model.IconMeta{
			ID:       util.Slugify(f.Name()), // Use the file name as ID
			Name:     f.Name(),
			FullPath: filepath.Join(dir, f.Name()),
			Size:     info.Size(),
		}
	}
	return nil
}

func resolveStatus(entries []os.DirEntry) string {
	hasIco := false
	hasIni := false

	for _, f := range entries {
		if f.IsDir() {
			continue
		}
		switch filepath.Ext(f.Name()) {
		case ".ico":
			hasIco = true
		case ".ini":
			if f.Name() == "desktop.ini" {
				hasIni = true
			}
		}

		// Early exit once both are found
		if hasIco && hasIni {
			return "ok"
		}
	}

	if hasIco {
		return "warn"
	}

	return "missing"
}

func sortedKeys[M ~map[string]V, V any](m M) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}
