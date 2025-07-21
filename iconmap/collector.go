package iconmap

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mrizkifadil26/medix/model"
	"github.com/mrizkifadil26/medix/utils"
)

func collectIcons(baseDir, source string, excludeDirs []string, dirMap map[string][]*model.IconEntry) error {
	baseIconMap := make(map[string]*model.IconEntry)
	var allIcos []string

	err := filepath.WalkDir(baseDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Printf("⚠️ Walk error in %s: %v", path, err)
			return nil
		}

		if d.IsDir() && isExcluded(filepath.Base(path), excludeDirs) {
			return filepath.SkipDir
		}

		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".ico") {
			allIcos = append(allIcos, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// First pass: add base icons
	for _, path := range allIcos {
		relDir, filename, baseName := getNames(baseDir, path)
		if filename != baseName {
			continue // skip variant for now
		}

		info, _ := os.Stat(path)
		entry := &model.IconEntry{
			ID:       utils.Slugify(filename),
			Name:     filename + ".ico",
			Size:     info.Size(),
			Source:   source,
			FullPath: path,
			Type:     "icon",
			Variants: []string{},
		}

		dirMap[relDir] = append(dirMap[relDir], entry)
		baseKey := filepath.Join(relDir, baseName)
		baseIconMap[baseKey] = entry

		// log.Printf("🎯 Base icon added: %s (relDir: %s)", entry.Name, relDir)
	}

	// Second pass: assign variants
	for _, path := range allIcos {
		relDir, filename, baseName := getNames(baseDir, path)
		if filename == baseName {
			continue // already handled
		}

		baseKey := filepath.Join(relDir, baseName)
		if entry, ok := baseIconMap[baseKey]; ok {
			entry.Variants = append(entry.Variants, path)
			// log.Printf("✅ Variant added: %s → base: %s", filename, entry.Name)
		} else {
			// no base found, add as fallback
			info, _ := os.Stat(path)
			fallback := &model.IconEntry{
				ID:       utils.Slugify(filename),
				Name:     filename + ".ico",
				Size:     info.Size(),
				Source:   source,
				FullPath: path,
				Type:     "icon",
				Variants: []string{},
			}
			dirMap[relDir] = append(dirMap[relDir], fallback)
			log.Printf("⚠️ No base found for variant %s, added as fallback", fallback.Name)
		}
	}

	return nil
}

func isExcluded(name string, excludeDirs []string) bool {
	return slices.Contains(excludeDirs, name)
}

func getNames(baseDir, fullPath string) (relDir, filename, baseName string) {
	relDir, _ = filepath.Rel(baseDir, filepath.Dir(fullPath))
	filename = strings.TrimSuffix(filepath.Base(fullPath), filepath.Ext(fullPath))
	baseName = stripAltSuffix(filename)

	return
}

// stripAltSuffix removes trailing variant markers like " (Alt)", " (Alt 2)", etc.
func stripAltSuffix(name string) string {
	name = strings.TrimSpace(name)
	if idx := strings.LastIndex(name, " (alt"); idx != -1 && strings.HasSuffix(name, ")") {
		return strings.TrimSpace(name[:idx])
	}

	return name
}
