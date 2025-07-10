package organize

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mrizkifadil26/medix/model"
)

func IndexScatteredIcons(baseDir, source string, excludeDirs []string) ([]model.IconEntry, error) {
	var result []model.IconEntry

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".ico") {
			continue
		}

		full := filepath.Join(baseDir, entry.Name())
		info, err := os.Stat(full)
		if err != nil {
			continue
		}

		result = append(result, model.IconEntry{
			ID:       "", // will be filled after slug mapping
			Name:     entry.Name(),
			Size:     info.Size(),
			FullPath: full,
			Source:   source,
			Type:     "icon",
		})
	}

	return result, nil
}
