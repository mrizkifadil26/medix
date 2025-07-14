package sync

import (
	"fmt"
	"time"

	"github.com/mrizkifadil26/medix/model"
	util "github.com/mrizkifadil26/medix/utils"
)

// GenerateUnusedIconReport builds a report based on icons not used in any media
func GenerateUnusedIconsReport(index *model.SyncedIconIndex) error {
	report := UnusedIconReport{
		GeneratedAt: time.Now(),
		Groups:      make(map[string][]UnusedIconEntry),
	}

	for _, group := range index.Data {
		for _, entry := range group.Items {
			if entry.UsedBy == nil {
				groupName := group.Name
				report.Groups[groupName] = append(report.Groups[groupName], UnusedIconEntry{
					Name:   entry.Name,
					Path:   entry.FullPath,
					Source: entry.Source,
				})
				report.Total++
			}
		}
	}

	if report.Total > 0 {
		fmt.Printf("🧾 Writing unused icons report (%d unused)...\n", report.Total)
		return util.WriteReport("unused-icons", report)
	}

	fmt.Println("✅ No unused icons found, skipping report.")
	return nil
}
