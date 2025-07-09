package sync

import (
	"fmt"
	"time"

	"github.com/mrizkifadil26/medix/model"
	"github.com/mrizkifadil26/medix/util"
)

// GenerateUnusedIconReport builds a report based on icons not used in any media
func GenerateUnusedIconsReport(index *model.SyncedIconIndex) error {
	var report UnusedIconReport
	report.GeneratedAt = time.Now()

	for _, group := range index.Data {
		for _, entry := range group.Items {
			if entry.UsedBy == nil {
				report.Icons = append(report.Icons, UnusedIconEntry{
					Name:     entry.Name,
					FullPath: entry.FullPath,
					Source:   entry.Source,
				})
			}
		}
	}

	report.Total = len(report.Icons)

	if report.Total > 0 {
		fmt.Printf("🧾 Writing unused icons report (%d unused)...\n", report.Total)
	} else {
		fmt.Println("✅ No unused icons found, skipping report.")
		return nil
	}

	return util.WriteReport("unused-icons", report)
}
