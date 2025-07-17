package iconmap

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mrizkifadil26/medix/utils"
)

func collectIcons(baseDir, source string, excludeDirs []string, dirMap map[string][]IconEntry) error {
	return filepath.WalkDir(baseDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Printf("⚠️ Walk error in %s: %v", path, err)
			return nil
		}

		if d.IsDir() && isExcluded(filepath.Base(path), excludeDirs) {
			return filepath.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(strings.ToLower(d.Name()), ".ico") {
			return nil
		}

		info, err := os.Stat(path)
		if err != nil {
			log.Printf("⚠️ Failed to stat %s: %v", path, err)
			return nil
		}

		relDir, err := filepath.Rel(baseDir, filepath.Dir(path))
		if err != nil || relDir == "." {
			return nil
		}

		dirMap[relDir] = append(dirMap[relDir], IconEntry{
			ID:       utils.Slugify(d.Name()),
			Name:     d.Name(),
			Size:     info.Size(),
			Source:   source,
			FullPath: path,
			Type:     "icon",
		})

		return nil
	})
}

func isExcluded(name string, excludeDirs []string) bool {
	return slices.Contains(excludeDirs, name)
}
