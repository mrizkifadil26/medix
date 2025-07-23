package scanner

import (
	"os"
	"path/filepath"

	"github.com/mrizkifadil26/medix/model"
	"github.com/mrizkifadil26/medix/utils"
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

func resolveIcon(dir string, entries []os.DirEntry) *model.IconRef {
	for _, f := range entries {
		if f.IsDir() || filepath.Ext(f.Name()) != ".ico" {
			continue
		}

		info, err := os.Stat(filepath.Join(dir, f.Name()))
		if err != nil {
			continue
		}

		name := f.Name()
		size := info.Size()
		return &model.IconRef{
			Name:     name,
			Slug:     utils.Slugify(name),
			FullPath: filepath.Join(dir, name),
			Size:     size,
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

func hasValidExt(path string, exts []string) bool {
	ext := filepath.Ext(path)
	for _, e := range exts {
		if e == ext {
			return true
		}
	}

	return false
}
