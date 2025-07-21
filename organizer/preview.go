package organizer

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/mrizkifadil26/medix/model"
)

func Preview(scan model.MediaOutput, sources []IconSource, targetDir string) OrganizeResult {
	iconMap := loadIcons(sources)
	result := OrganizeResult{
		Type:        scan.Source,
		GeneratedAt: time.Now(),
	}

	for _, item := range scan.Items {
		if item.Status == "missing" {
			continue
		}

		slug := item.Icon.ID
		icons, ok := iconMap[slug]

		if !ok {
			continue // no icon found
		}

		for _, icon := range icons {
			if icon.Size == item.Icon.Size {
				targetPath := filepath.Join(targetDir, item.Group, filepath.Base(icon.FullPath))

				fmt.Printf("✅ Matched: %-40s → %-20s [size: %d]\n", item.Name, filepath.Base(icon.FullPath), icon.Size)

				result.Changes = append(result.Changes, ChangeItem{
					Action: "move",
					Source: icon.FullPath,
					Target: targetPath,
					Reason: "slug+size match",
					Item:   item.BaseEntry,
				})

				break
			}
		}
	}

	return result
}
