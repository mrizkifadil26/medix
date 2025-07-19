package organizer

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mrizkifadil26/medix/model"
	"github.com/mrizkifadil26/medix/utils"
)

func Preview(scan model.MediaOutput, sources []string) OrganizeResult {
	iconMap := loadIcons(sources)
	result := OrganizeResult{
		Type:        scan.Source,
		GeneratedAt: time.Now(),
	}

	for _, item := range scan.Items {
		if item.Status != "missing" {
			continue
		}

		slug := utils.Slugify(item.Name)
		iconPath, ok := iconMap[slug]
		if !ok {
			continue // no icon found
		}

		result.Changes = append(result.Changes, ChangeItem{
			Action: "move",
			Source: iconPath,
			Target: filepath.Join(item.Path, "folder.ico"),
			Reason: "slug match",
			Item:   item.BaseEntry,
		})
	}

	return result
}

// loadIcons builds a map[slug]filepath from source directories
func loadIcons(sources []string) map[string]string {
	iconMap := make(map[string]string)

	for _, src := range sources {
		entries, err := os.ReadDir(src)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			if !strings.HasSuffix(strings.ToLower(entry.Name()), ".ico") {
				continue
			}

			name := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			slug := utils.Slugify(name)

			// only set if not already mapped
			if _, exists := iconMap[slug]; !exists {
				iconMap[slug] = filepath.Join(src, entry.Name())
			}
		}
	}

	return iconMap
}
