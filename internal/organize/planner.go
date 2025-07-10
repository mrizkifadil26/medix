package organize

import (
	"os"
	"path/filepath"

	"github.com/mrizkifadil26/medix/model"
)

type MovePlan struct {
	SourcePath string `json:"source_path"`
	TargetPath string `json:"target_path"`
	MatchedID  string `json:"matched_id"`
	Group      string `json:"group"`
	Type       string `json:"type"` // movie or tv_show
	Matched    bool   `json:"matched"`
	Duplicate  bool   `json:"duplicate"`
}

func BuildMovePlan(icons []model.IconEntry, slugMap MediaSlugMap, outputBase string) []MovePlan {
	var plan []MovePlan

	for _, icon := range icons {
		entry, ok := MatchIconToEntry(icon, slugMap)
		if !ok {
			continue
		}

		group := entry.Group
		targetDir := filepath.Join(outputBase, group)
		targetPath := filepath.Join(targetDir, icon.Name)

		_, err := os.Stat(targetPath)
		isDuplicate := err == nil

		plan = append(plan, MovePlan{
			SourcePath: icon.FullPath,
			TargetPath: targetPath,
			MatchedID:  entry.Icon.ID,
			Group:      group,
			Type:       entry.Type,
			Matched:    true,
			Duplicate:  isDuplicate,
		})
	}

	return plan
}
