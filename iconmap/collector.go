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

func ScanIconDirectory(baseDir, source string, excludeDirs []string) (map[string][]model.IconEntry, error) {
	groupedIcons := make(map[string][]*model.IconEntry)
	baseIconMap := make(map[string]*model.IconEntry)
	iconPaths := []string{}

	// First pass: collect all .ico file paths
	err := filepath.WalkDir(baseDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Printf("⚠️ Walk error in %s: %v", path, err)
			return nil
		}

		if d.IsDir() && isExcluded(filepath.Base(path), excludeDirs) {
			return filepath.SkipDir
		}

		if !d.IsDir() && strings.HasSuffix(strings.ToLower(d.Name()), ".ico") {
			iconPaths = append(iconPaths, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Second pass: register base icons
	for _, path := range iconPaths {
		relDir, filename, baseName := getNames(baseDir, path)

		// Skip if it's a variant
		if filename != baseName {
			continue
		}

		info, _ := os.Stat(path)
		entry := &model.IconEntry{
			ID:       utils.Slugify(filename),
			Name:     filename + ".ico",
			Size:     info.Size(),
			Source:   source,
			Path:     path,
			Type:     "icon",
			Variants: []string{},
		}

		groupedIcons[relDir] = append(groupedIcons[relDir], entry)
		baseKey := filepath.Join(relDir, baseName)
		baseIconMap[baseKey] = entry

		// log.Printf("🎯 Base icon added: %s (relDir: %s)", entry.Name, relDir)
	}

	// Third pass: assign variants or fallback if no base found
	for _, path := range iconPaths {
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
				Path:     path,
				Type:     "icon",
				Variants: []string{},
			}

			groupedIcons[relDir] = append(groupedIcons[relDir], fallback)
			log.Printf("⚠️ No base found for variant %s, added as fallback", fallback.Name)
		}
	}

	finalGrouped := make(map[string][]model.IconEntry)
	for group, entries := range groupedIcons {
		finalGrouped[group] = toValueSlice(entries)
	}

	return finalGrouped, nil
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

func toValueSlice(entries []*model.IconEntry) []model.IconEntry {
	result := make([]model.IconEntry, len(entries))
	for i, entry := range entries {
		result[i] = *entry
	}
	return result
}
