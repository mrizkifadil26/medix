package organizer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mrizkifadil26/medix/model"
	"github.com/mrizkifadil26/medix/utils"
)

func Preview(scan model.MediaOutput, sources []IconSource, targetDir string) OrganizeResult {
	iconMap := loadIcons(sources)
	result := OrganizeResult{
		Type:        scan.Source,
		GeneratedAt: time.Now(),
	}

	for _, item := range scan.Items {
		// if item.Status != "missing" {
		// 	continue
		// }

		slug := utils.Slugify(item.Name)
		iconPath, ok := iconMap[slug]
		if !ok {
			continue // no icon found
		}

		result.Changes = append(result.Changes, ChangeItem{
			Action: "copy",
			Source: iconPath,
			Target: filepath.Join(targetDir, item.Group, filepath.Base(iconPath)),
			Reason: "slug match",
			Item:   item.BaseEntry,
		})
	}

	return result
}

// loadIcons builds a map[slug]filepath from source directories
func loadIcons(sources []IconSource) map[string]string {
	iconMap := make(map[string]string)

	for _, src := range sources {
		entries, err := os.ReadDir(src.Path)
		if err != nil {
			fmt.Printf("⚠️ Failed to read icon source %s (%s): %v\n", src.Source, src.Path, err)
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
				iconMap[slug] = filepath.Join(src.Path, entry.Name())
			}
		}
	}

	return iconMap
}
